package audio

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProcessAudioRequest(t *testing.T) {
	// Create a new Echo instance to simulate requests
	e := echo.New()

	// Define a test case
	testCases := []struct {
		name        string
		requestBody string
		expectedErr error
	}{
		{
			name:        "Valid Request",
			requestBody: `{"file":"file1","model":"whisper-1","prompt":"example prompt","response_format":"json","temperature":0.7,"language":"en"}`,
			expectedErr: nil,
		},
		{
			name:        "Missing Required Field",
			requestBody: `{"model":"whisper-1","prompt":"example prompt","response_format":"json","temperature":0.7,"language":"en"}`,
			expectedErr: echo.NewHTTPError(http.StatusBadRequest, "Required parameters are not set"),
		},
		// Add more test cases to cover different scenarios.
	}
	// Set the OpenAI API key for testing
	viper.Set("openAI.apiKey", "YOUR_RANDOM_API_KEY")

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// Create a new response recorder to capture the response
			rec := httptest.NewRecorder()

			// Perform the request/response cycle
			c := e.NewContext(req, rec)
			err := processAudioRequest(c, "https://example.com")

			// Check the error
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
		})
	}
}

// MockContext creates a mock echo.Context for testing purposes
func MockContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec)
}

func TestProcessAudioRequest_Unauthorized(t *testing.T) {
	// Create a mock context
	c := MockContext()

	// Clear the apiKey so that it will be empty
	viper.Set("openAI.apiKey", "")

	// Set the URL to an OpenAI API URL
	url := "https://api.openai.com/v1/audio/transcriptions"

	// Call the processAudioRequest function
	err := processAudioRequest(c, url)

	// Check if the error is the expected unauthorized error
	if err == nil || !strings.Contains(err.Error(), "OpenAI API key not found") {
		t.Errorf("Expected unauthorized error, got: %v", err)
	}
}

func TestGetParams_Success(t *testing.T) {
	// Create a mock context
	c := MockContext()

	// Set the request body with valid RequestBody data
	reqBody := `{
		"file": "example.wav",
		"model": "whisper-1",
		"prompt": "This is a prompt",
		"response_format": "json",
		"temperature": 0.8,
		"language": "en"
	}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c = echo.New().NewContext(req, rec)

	// Call the getRequestBody function
	params, err := getRequestBody(c)

	// Check if there is no error
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check if the returned RequestBody match the expected values
	expectedParams := RequestBody{
		File:           "example.wav",
		Model:          "whisper-1",
		Prompt:         "This is a prompt",
		ResponseFormat: "json",
		Temperature:    0.8,
		Language:       "en",
	}

	assert.ObjectsAreEqualValues(params, expectedParams)
}

func TestGetParams_UnmarshallParametersFail(t *testing.T) {
	// Create a mock context
	c := MockContext()

	// Set the request body with missing required parameters
	reqBody := `{
		"file": "example.wav",
		"model": "whisper-1",
		"prompt": "This is a prompt",
		"response_format": "json",
		"temperature": "ERROR",
		"language": "en",
		"extra_param": "extra"
	}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c = echo.New().NewContext(req, rec)

	// Call the getRequestBody function
	params, err := getRequestBody(c)

	// Check if the error is the expected error indicating missing required parameters
	assert.Error(t, err)
	assert.ErrorContains(t, err, "Unmarshal type error")
	assert.Empty(t, params)
}

func TestMakeRequest(t *testing.T) {
	// Prepare the test data
	url := "https://api.openai.com/v1/audio/transcriptions"
	params := RequestBody{
		File:           "example.wav",
		Model:          "whisper-1",
		Prompt:         "This is a prompt",
		ResponseFormat: "json",
		Temperature:    0.8,
		Language:       "en",
	}

	// Mock the API key
	apiKey := "your_openai_api_key"

	// Call the makeRequest function
	req, err := makeRequest(&params, url, apiKey)

	// Check if there is no error
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check if the request attributes are correct
	if req.Method != http.MethodPost {
		t.Errorf("Expected POST method, got: %s", req.Method)
	}

	if req.URL.String() != url {
		t.Errorf("Expected URL %s, got: %s", url, req.URL.String())
	}

	if req.Header.Get("Authorization") != "Bearer "+apiKey {
		t.Errorf("Expected Authorization header with API key, got: %s", req.Header.Get("Authorization"))
	}

	if req.Header.Get("Content-Type") != "multipart/form-data" {
		t.Errorf("Expected Content-Type header application/json, got: %s", req.Header.Get("Content-Type"))
	}

	// Read and compare the request body
	expectedBody, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Error marshaling expected body: %v", err)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}

	if !bytes.Equal(body, expectedBody) {
		t.Errorf("Expected request body:\n%s\ngot:\n%s", expectedBody, body)
	}
}
