package api

import (
	"bytes"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
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
			err := processAudioRequest(c, "https://example.com", nil)

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
