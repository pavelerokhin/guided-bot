package api

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProcessCompletionsRequest_Success(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	reqBody := `{"model": "gpt-3.5-turbo", "prompt": "Hello, ChatGPT!"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockResponse := `{"id": "chatcmpl-123", "object": "chat.completion", "created": 1677652288, "choices": [{"index": 0, "message": {"role": "assistant", "content": "\\n\\nHello there, how may I assist you today?"}, "finish_reason": "stop"}], "usage": {"prompt_tokens": 9, "completion_tokens": 12, "total_tokens": 21}}`

	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "multipart/form-data", r.Header.Get("Content-Type"))

		// Verify request body
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)
		assert.JSONEq(t, reqBody, buf.String())

		// Respond with the JSON response
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockResponse))
		assert.Nil(t, err)
	}))
	defer testServer.Close()

	// Use the test server URL in the test case
	err := processCompletionsRequest(c, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestProcessCompletionsRequest_Unauthorized(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set an empty apiKey and a url that starts with "https://api.openai.com/"
	viper.Set("openAI.apiKey", "")
	url := "https://api.openai.com/some/endpoint"

	// Call the function being tested
	err := processCompletionsRequest(c, url)

	// Assert that the response is an HTTP 401 (Unauthorized) error
	assert.Error(t, err)
	// Assert that the response body contains the expected error message
	expectedErrorMessage := "code=401, message=OpenAI API key not found"
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestProcessCompletionsRequest_SendRequestError(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	reqBody := `ERROR`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a test server that returns an error response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with the JSON response
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("mocked API error"))
	}))
	defer testServer.Close()

	// Call the function being tested, using the test server URL for the API call
	err := processCompletionsRequest(c, testServer.URL)

	// Assert that there was an error during the API request
	assert.Error(t, err)
}
