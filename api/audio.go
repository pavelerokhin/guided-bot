package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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
	return processAudioRequest(c, "https://api.openai.com/v1/audio/transcriptions", nil)
}

func CreateTranslation(c echo.Context) error {
	return processAudioRequest(c, "https://api.openai.com/v1/audio/translations", nil)
}

func processAudioRequest(c echo.Context, url string, httpClient *http.Client) error {
	openAPIKey := viper.GetString("openAI.apiKey")
	if openAPIKey == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	params := AudioParameters{}
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if params.Model == "" || params.File == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Required parameters are not set")
	}

	// Convert the fields to JSON format
	data, err := json.Marshal(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", "multipart/form-data")

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}
