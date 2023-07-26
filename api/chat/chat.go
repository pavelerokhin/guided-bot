package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"OpenAI-api/api"
)

const (
	CompletionUrl = "https://api.openai.com/v1/chat/completions"
)

func CreateChatCompletion(c echo.Context) error {
	return processChatRequest(c, CompletionUrl)
}

func processChatRequest(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	requestBody, err := getRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := makeRequest(requestBody, url, apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := api.SendRequest(nil, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getRequestBody(c echo.Context) (*RequestBody, error) {
	requestBody := RequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Messages == nil {
		return nil, errors.New("required parameters are not set (required: Model, Messages)")
	}

	return &requestBody, nil
}

func makeRequest(params *RequestBody, url, apiKey string) (*http.Request, error) {
	// Convert the fields to JSON format
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Create a new POST request, set authorization header (if requested) and content type for audio
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	return req, nil
}
