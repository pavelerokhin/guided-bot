package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"OpenAI-api/api/request"
)

func TestCreateGetRequestBody_ValidRequest(t *testing.T) {
	// Prepare a valid request body
	jsonStr := `{
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageCreateRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "prompt", requestBody.Prompt)
}

func TestCreateGetRequestBody_InvalidRequest_MissingModel(t *testing.T) {
	// Prepare an invalid request body without the required "model" field
	jsonStr := `{
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageCreateRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Prompt)")
}

func TestCreateGetRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := request.GetImageCreateRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}

func TestCreateProcessChatRequest_Success(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	reqBody := `{"prompt": "Hello, ChatGPT!"}`
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
	err := processImageCreate(c, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestCreateProcessChatRequest_Unauthorized(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set an empty apiKey and a url that starts with "https://api.openai.com/"
	viper.Set("openAI.apiKey", "")
	url := "https://api.openai.com/some/endpoint"

	// Call the function being tested
	err := processImageCreate(c, url)

	// Assert that the response is an HTTP 401 (Unauthorized) error
	assert.Error(t, err)
	// Assert that the response body contains the expected error message
	expectedErrorMessage := "code=401, message=OpenAI API key not found"
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestCreateProcessChatRequest_SendRequestError(t *testing.T) {
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
	err := processImageCreate(c, testServer.URL)

	// Assert that there was an error during the API request
	assert.Error(t, err)
}

func TestEditGetRequestBody_ValidRequest(t *testing.T) {
	// Prepare a valid request body
	jsonStr := `{
		"image": "some_image",
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageEditRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "prompt", requestBody.Prompt)
}

func TestEditGetRequestBody_InvalidRequest_MissingImage(t *testing.T) {
	// Prepare an invalid request body without the required "model" field
	jsonStr := `{
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageEditRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Image, Prompt)")
}

func TestEditGetRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := request.GetImageEditRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}

func TestEditProcessChatRequest_Success(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	reqBody := `{"image": "some_image", "prompt": "Hello, ChatGPT!"}`
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
	err := processImageEdit(c, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestEditProcessChatRequest_Unauthorized(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set an empty apiKey and a url that starts with "https://api.openai.com/"
	viper.Set("openAI.apiKey", "")
	url := "https://api.openai.com/some/endpoint"

	// Call the function being tested
	err := processImageEdit(c, url)

	// Assert that the response is an HTTP 401 (Unauthorized) error
	assert.Error(t, err)
	// Assert that the response body contains the expected error message
	expectedErrorMessage := "code=401, message=OpenAI API key not found"
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestEditProcessChatRequest_SendRequestError(t *testing.T) {
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
	err := processImageEdit(c, testServer.URL)

	// Assert that there was an error during the API request
	assert.Error(t, err)
}

func TestVariationGetRequestBody_ValidRequest(t *testing.T) {
	// Prepare a valid request body
	jsonStr := `{
		"image": "some_image",
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageVariationRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "some_image", requestBody.Image)
}

func TestVariationGetRequestBody_InvalidRequest_MissingImage(t *testing.T) {
	// Prepare an invalid request body without the required "model" field
	jsonStr := `{
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := request.GetImageVariationRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Image)")
}

func TestVariationGetRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := request.GetImageVariationRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}

func TestVariationProcessChatRequest_Success(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	reqBody := `{"image": "some_image", "prompt": "Hello, ChatGPT!"}`
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
	err := processImageEdit(c, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestVariationProcessChatRequest_Unauthorized(t *testing.T) {
	// Set up the test Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set an empty apiKey and a url that starts with "https://api.openai.com/"
	viper.Set("openAI.apiKey", "")
	url := "https://api.openai.com/some/endpoint"

	// Call the function being tested
	err := processImageVariation(c, url)

	// Assert that the response is an HTTP 401 (Unauthorized) error
	assert.Error(t, err)
	// Assert that the response body contains the expected error message
	expectedErrorMessage := "code=401, message=OpenAI API key not found"
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestVariationProcessChatRequest_SendRequestError(t *testing.T) {
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
	err := processImageVariation(c, testServer.URL)

	// Assert that there was an error during the API request
	assert.Error(t, err)
}
