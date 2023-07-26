package audio

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

const (
	transcriptionUrl = "https://api.openai.com/v1/audio/transcriptions"
	translationUrl   = "https://api.openai.com/v1/audio/translations"
)

func CreateTranscription(c echo.Context) error {
	return processAudioRequest(c, transcriptionUrl)
}

func CreateTranslation(c echo.Context) error {
	return processAudioRequest(c, translationUrl)
}

func processAudioRequest(c echo.Context, url string) error {
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

	body, err := sendRequest(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getRequestBody(c echo.Context) (*RequestBody, error) {
	params := RequestBody{}

	if err := c.Bind(&params); err != nil {
		return nil, err
	}

	if params.Model == "" || params.File == "" {
		return &params, errors.New("required parameters are not set (required: Model, File)")
	}

	return &params, nil
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
