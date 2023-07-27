package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"OpenAI-api/api/model"
)

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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
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

	req.Header.Set("Content-Type", "multipart/form-data")

	return req, nil
}
