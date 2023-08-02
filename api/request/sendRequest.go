package request

import (
	"fmt"
	"io"
	"net/http"
)

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
