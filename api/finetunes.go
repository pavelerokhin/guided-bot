package api

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
)

// FineTuneRequest represents the JSON request body for fine-tuning
type FineTuneRequest struct {
	TrainingFile string `json:"training_file"`
	// Add other fields as needed based on the request body
}

func CreateFineTune(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	reqBody := FineTuneRequest{
		TrainingFile: "file-XGinujblHPwGLSztz8cPS8XY", // Replace with the ID of the uploaded training file
		// Add other fields as needed based on the request body
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to marshal request data")
	}

	// Create the POST request
	url := "https://api.openai.com/v1/fine-tunes"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request to create the fine-tuning job
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create fine-tuning job")
	}
	defer resp.Body.Close()

	// Read the response body
	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the response data as JSON
	return c.JSON(http.StatusOK, responseData)
}
