package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"OpenAI-api/api/model"
)

type ClientMock struct {
}

func (c ClientMock) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("mock HTTP client error")
}

func TestMakeRequest_NoAPIKey(t *testing.T) {
	// Prepare valid parameters for the function but without an API key
	requestBody := &model.CompletionsRequestBody{
		Model:  "some_model",
		Prompt: "prompt",
	}
	url := "https://api.example.com/completions"
	apiKey := ""

	// Call the function
	req, err := MakeRequest(requestBody, url, apiKey)

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

func TestSendRequest_SuccessfulResponse(t *testing.T) {
	// Prepare a mock HTTP server
	mockResponse := `{"result": "success"}`
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, mockResponse)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Create a request to the mock server
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	// Call the function
	response, err := SendRequest(nil, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []byte(mockResponse), response)
}

func TestSendRequest_ErrorOnRequest(t *testing.T) {
	// Create a request with an invalid URL to trigger an error
	req, err := http.NewRequest("GET", "invalid-url", nil)
	assert.NoError(t, err)

	// Call the function
	response, err := SendRequest(nil, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestSendRequest_ErrorOnReadResponse(t *testing.T) {
	// Prepare a mock HTTP server that returns an empty response
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Create a request to the mock server
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	// Call the function
	response, err := SendRequest(nil, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, strings.Contains(err.Error(), "unexpected status code"))
}

func TestSendRequest_HTTPClientError(t *testing.T) {
	// Create a request to the mock server
	req, err := http.NewRequest("GET", "https://api.example.com", nil)
	assert.NoError(t, err)

	// Create a custom HTTP client (the mock client) for testing
	mockClient := ClientMock{}

	// Call the function with the custom HTTP client
	response, err := SendRequest(mockClient, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.EqualError(t, err, "mock HTTP client error")
}

func TestSendRequest_ReadResponseBody(t *testing.T) {
	// Create a mock HTTP server
	mockResponse := `{"result": "success"}`
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, mockResponse)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Create a request to the mock server
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	// Call the function
	response, err := SendRequest(http.DefaultClient, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []byte(mockResponse), response)
}

func TestSendRequest_ReadResponseBodyError(t *testing.T) {
	// Create a request to the mock server (not used in this test)
	req, err := http.NewRequest("GET", "https://api.example.com", nil)
	assert.NoError(t, err)

	// Create a custom HTTP client (the mock client) for testing
	mockClient := &ClientMock{}

	// Call the function with the custom HTTP client
	response, err := SendRequest(mockClient, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.EqualError(t, err, "mock HTTP client error")
}

func TestMakeRequest_ChatValidParams(t *testing.T) {
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
	req, err := MakeRequest(requestBody, url, apiKey)

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

func TestMakeRequest_CompletionsValidParams(t *testing.T) {
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
	req, err := MakeRequest(requestBody, url, apiKey)

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

func TestMakeRequest_ImageCreateValidParams(t *testing.T) {
	// Prepare valid parameters for the function
	requestBody := &model.ImageCreateRequestBody{
		Prompt: "some_model",
	}
	url := "https://api.example.com/chat"
	apiKey := "YOUR_API_KEY"

	// Call the function
	req, err := MakeRequest(requestBody, url, apiKey)

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

func TestMakeRequest_ImageEditValidParams(t *testing.T) {
	// Prepare valid parameters for the function
	requestBody := &model.ImageEditRequestBody{
		Image:  "some_image",
		Prompt: "some prompt",
	}
	url := "https://api.example.com/chat"
	apiKey := "YOUR_API_KEY"

	// Call the function
	req, err := MakeRequest(requestBody, url, apiKey)

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

func TestMakeRequest_ImageVariateValidParams(t *testing.T) {
	// Prepare valid parameters for the function
	requestBody := &model.ImageVariateRequestBody{
		Image: "some_image",
	}
	url := "https://api.example.com/chat"
	apiKey := "YOUR_API_KEY"

	// Call the function
	req, err := MakeRequest(requestBody, url, apiKey)

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
