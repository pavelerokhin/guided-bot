package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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

// FineTuneJob represents the JSON response for a fine-tuning job
type FineTuneJob struct {
	ID              string                   `json:"id"`
	Object          string                   `json:"object"`
	Model           string                   `json:"model"`
	CreatedAt       int64                    `json:"created_at"`
	Events          []map[string]interface{} `json:"events"`
	FineTunedModel  string                   `json:"fine_tuned_model"`
	Hyperparams     map[string]interface{}   `json:"hyperparams"`
	OrganizationID  string                   `json:"organization_id"`
	ResultFiles     []interface{}            `json:"result_files"`
	Status          string                   `json:"status"`
	ValidationFiles []interface{}            `json:"validation_files"`
	TrainingFiles   []interface{}            `json:"training_files"`
	UpdatedAt       int64                    `json:"updated_at"`
}

func ListFineTunes(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Create the GET request
	url := "https://api.openai.com/v1/fine-tunes"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to list fine-tuning jobs
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to list fine-tuning jobs")
	}
	defer resp.Body.Close()

	// Read the response body
	var responseData map[string][]FineTuneJob
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the response data as JSON
	return c.JSON(http.StatusOK, responseData)
}

func RetrieveFineTunes(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Get the fine_tune_id from the path parameters
	fineTuneID := c.Param("fine_tune_id")

	// Create the GET request
	url := fmt.Sprintf("https://api.openai.com/v1/fine-tunes/%s", fineTuneID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to retrieve information about the fine-tune job
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve fine-tune job")
	}
	defer resp.Body.Close()

	// Read the response body
	var responseData FineTuneJob
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the response data as JSON
	return c.JSON(http.StatusOK, responseData)
}

func CancelFineTunes(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Get the fine_tune_id from the path parameters
	fineTuneID := c.Param("fine_tune_id")

	// Create the POST request
	url := fmt.Sprintf("https://api.openai.com/v1/fine-tunes/%s/cancel", fineTuneID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to cancel the fine-tune job
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to cancel fine-tune job")
	}
	defer resp.Body.Close()

	// Read the response body
	var responseData FineTuneJob
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the response data as JSON
	return c.JSON(http.StatusOK, responseData)
}

// FineTuneEvent represents a fine-tune event
type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func ListFineTuneEvents(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Get the fine_tune_id from the path parameters
	fineTuneID := c.Param("fine_tune_id")

	// Create the GET request
	url := fmt.Sprintf("https://api.openai.com/v1/fine-tunes/%s/events", fineTuneID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to retrieve fine-tune events
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve fine-tune events")
	}
	defer resp.Body.Close()

	// Read the response body
	var responseData struct {
		Object string          `json:"object"`
		Data   []FineTuneEvent `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Return the response data as JSON
	return c.JSON(http.StatusOK, responseData)
}

func deleteFineTunedModelHandler(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Get the model name from the path parameters
	modelName := c.Param("model")

	// Create the DELETE request
	url := fmt.Sprintf("https://api.openai.com/v1/models/%s", modelName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to delete the fine-tuned model
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete the fine-tuned model")
	}
	defer resp.Body.Close()

	// Return the response as JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":      modelName,
		"object":  "model",
		"deleted": true,
	})
}
