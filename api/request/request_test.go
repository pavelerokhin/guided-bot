package request

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
}

func (c ClientMock) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("mock HTTP client error")
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
