package request

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetChatRequestBody_ValidRequest(t *testing.T) {
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
	requestBody, err := GetChatRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "some_model", requestBody.Model)
	assert.Len(t, requestBody.Messages, 1)
	assert.Equal(t, "user", requestBody.Messages[0].Role)
	assert.Equal(t, "Hello, ChatGPT!", requestBody.Messages[0].Content)
}

func TestGetChatRequestBody_InvalidRequest_MissingModel(t *testing.T) {
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
	requestBody, err := GetChatRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Messages)")
}

func TestGetChatRequestBody_InvalidRequest_MissingMessages(t *testing.T) {
	// Prepare an invalid request body without the required "messages" field
	jsonStr := `{
		"model": "some_model"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := GetChatRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Messages)")
}

func TestGetChatRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := GetChatRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}

func TestCompletionsGetRequestBody_ValidRequest(t *testing.T) {
	// Prepare a valid request body
	jsonStr := `{
		"model": "some_model",
		"prompt": "prompt"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := GetCompletionsRequestBody(c)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, requestBody)
	assert.Equal(t, "some_model", requestBody.Model)
	assert.Equal(t, "prompt", requestBody.Prompt)
}

func TestCompletionsGetRequestBody_InvalidRequest_MissingModel(t *testing.T) {
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
	requestBody, err := GetCompletionsRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Prompt)")
}

func TestCompletionsGetRequestBody_InvalidRequest_MissingMessages(t *testing.T) {
	// Prepare an invalid request body without the required "messages" field
	jsonStr := `{
		"model": "some_model"
	}`
	req := httptest.NewRequest("POST", "/some-endpoint", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Call the function
	requestBody, err := GetCompletionsRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "required parameters are not set (required: Model, Prompt)")
}

func TestCompletionsGetRequestBody_BindError(t *testing.T) {
	// Create a new mock Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid request body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the function with the mock context
	requestBody, err := GetCompletionsRequestBody(c)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, requestBody)
	assert.EqualError(t, err, "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value")
}
