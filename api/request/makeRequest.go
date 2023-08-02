package request

import (
	"bytes"
	"encoding/json"
	"net/http"

	"OpenAI-api/api/model"
)

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
