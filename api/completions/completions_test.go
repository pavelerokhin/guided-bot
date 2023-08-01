package completions

import (
	"OpenAI-api/api/model"
	"OpenAI-api/api/request"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequestBody_ValidRequest(t *testing.T) {
	// Prepare a valid request body
	jsonStr := `{
		"model": "some_model",
		"messages": [{
			"role": "user",
			"content": "Hello, ChatGPT!"
		}]
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := getRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "some_model", requestBody.Model)
	assert.Len(t, requestBody.Prompt, 1)
}

func TestGetRequestBody_InvalidRequest_MissingModel(t *testing.T) {
	// Prepare an invalid request body without the required "model" field
	jsonStr := `{
		"messages": [{
			"role": "user",
			"content": "Hello, ChatGPT!"
		}]
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := getRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Messages)")
}

func TestGetRequestBody_InvalidRequest_MissingMessages(t *testing.T) {
	// Prepare an invalid request body without the required "messages" field
	jsonStr := `{
		"model": "some_model"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := getRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Messages)")
}

func TestMakeRequest_ValidParams(t *testing.T) {
	// Prepare valid parameters for the function
	requestBody := &model.ChatRequestBody{
		Model: "some_model",
		Messages: []model.Message{
			{
				Role:    "user",
				Content: "Hello, ChatGPT!",
			},
		},
	}
	url := "https://api.example.com/chat"
	apiKey := "YOUR_API_KEY"

	// Call the function
	req, err := request.MakeRequest(requestBody, url, apiKey)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, url, req.URL.String())

	// Check the authorization header
	assert.Equal(t, "Bearer "+apiKey, req.Header.Get("Authorization"))

	// Check the content type header
	assert.Equal(t, "multipart/form-data", req.Header.Get("Content-Type"))

	// Check the request body data
	expectedData, _ := json.Marshal(requestBody)
	// Read the request body and compare its contents to the expectedData
	bodyData, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, bodyData)
}

func TestMakeRequest_NoAPIKey(t *testing.T) {
	// Prepare valid parameters for the function but without an API key
	requestBody := &model.ChatRequestBody{
		Model: "some_model",
		Messages: []model.Message{
			{
				Role:    "user",
				Content: "Hello, ChatGPT!",
			},
		},
	}
	url := "https://api.example.com/chat"
	apiKey := ""

	// Call the function
	req, err := request.MakeRequest(requestBody, url, apiKey)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, url, req.URL.String())

	// Check that there is no authorization header
	assert.Empty(t, req.Header.Get("Authorization"))

	// Check the content type header
	assert.Equal(t, "multipart/form-data", req.Header.Get("Content-Type"))

	// Check the request body data
	expectedData, _ := json.Marshal(requestBody)
	// Read the request body and compare its contents to the expectedData
	bodyData, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, bodyData)
}

func TestGetRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := getRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}

func TestProcessChatRequest_Success(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	reqBody := `{"model": "gpt-3.5-turbo", "messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}]}`
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
	err := processChatRequest(c, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestProcessChatRequest_Unauthorized(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set an empty apiKey and a url that starts with "https://api.openai.com/"
	viper.Set("openAI.apiKey", "")
	url := "https://api.openai.com/some/endpoint"

	// Call the function being tested
	err := processChatRequest(c, url)

	// Assert that the response is an HTTP 401 (Unauthorized) error
	assert.Error(t, err)
	// Assert that the response body contains the expected error message
	expectedErrorMessage := "code=401, message=OpenAI API key not found"
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestProcessChatRequest_SendRequestError(t *testing.T) {
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
	err := processChatRequest(c, testServer.URL)

	// Assert that there was an error during the API request
	assert.Error(t, err)
}
