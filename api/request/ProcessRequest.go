package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"OpenAI-api/api/model"
)

func ProcessRequest[T model.RequestBody](url string, requestBody *T) ([]byte, error) {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	req, err := MakeRequest(requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func MakeRequest[T model.RequestBody](params *T, url, apiKey string) (*http.Request, error) {
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

	req.Header.Set("Content-Type", "application/json") // TODO: check multipart/form-data

	return req, nil
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func SendRequest(client HttpClient, req *http.Request) ([]byte, error) {
	if client == nil {
		client = http.DefaultClient
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	// Check if the response has a non-200 status code
	if resp.StatusCode != http.StatusOK {
		if resp != nil {
			// decode resp to get the error message
			var errorResponse model.Error
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		}
		return nil, fmt.Errorf("unexpected status code: %d. ", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
