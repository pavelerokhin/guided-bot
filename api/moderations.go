package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type ModerationInput struct {
	Input string `json:"input"`
}

type ModerationResult struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged        bool               `json:"flagged"`
		Categories     map[string]bool    `json:"categories"`
		CategoryScores map[string]float64 `json:"category_scores"`
	} `json:"results"`
}

func CreateModeration(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	// Parse the input text from the request body
	var input ModerationInput
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}

	// Create the request payload
	payload, err := json.Marshal(input)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request payload")
	}

	// Create the POST request
	url := "https://api.openai.com/v1/moderations"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Content-Type and Authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to perform content moderation
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to perform content moderation")
	}
	defer resp.Body.Close()

	// Parse the response JSON
	var result ModerationResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse response")
	}

	// Return the response as JSON
	return c.JSON(http.StatusOK, result)
}
