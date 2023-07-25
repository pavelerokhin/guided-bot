package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type AudioParameters struct {
	File           string  `json:"file"`
	Model          string  `json:"model"`  // ID of the model to use. Only whisper-1 is currently available.
	Prompt         string  `json:"prompt"` // If set should match audio language
	ResponseFormat string  `json:"response_format"`
	Temperature    float32 `json:"temperature"`
	Language       string  `json:"language"` // The language of the input audio. Supplying the input language in ISO-639-1 format will improve accuracy and latency.
}

func CreateTranscription(c echo.Context) error {
	return processAudioRequest(c, "https://api.openai.com/v1/audio/transcriptions")
}

func CreateTranslation(c echo.Context) error {
	return processAudioRequest(c, "https://api.openai.com/v1/audio/translations")
}

func processAudioRequest(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	params, err := getParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := makeRequest(params, url, apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := sendRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getParams(c echo.Context) (*AudioParameters, error) {
	params := AudioParameters{}

	if err := c.Bind(&params); err != nil {
		return nil, err
	}

	if params.Model == "" || params.File == "" {
		return &params, errors.New("Required parameters are not set")
	}

	return &params, nil
}

func makeRequest(params *AudioParameters, url, openAPIKey string) (*http.Request, error) {
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

	if openAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+openAPIKey)
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	return req, nil
}

func sendRequest(req *http.Request) ([]byte, error) {
	httpClient := http.DefaultClient

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return body, nil
}
