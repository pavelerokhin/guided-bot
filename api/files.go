package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	FileName  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

func ListFiles(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	url := "https://api.openai.com/v1/files"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var files []File
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   files,
		"object": "list",
	})
}

type FileUploadResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	FileName  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// TODO: Make this work
func UploadFile(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	purpose := c.FormValue("purpose")
	_, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Prepare the request body as a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add other form fields
	_ = writer.WriteField("purpose", purpose)

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create the HTTP request
	url := "https://api.openai.com/v1/files"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Set the Content-Type header for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteFile(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")
	fileID := c.Param("file_id")

	// Create the DELETE request
	url := fmt.Sprintf("https://api.openai.com/v1/files/%s", fileID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to delete the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete file")
	}
	defer resp.Body.Close()

	// Check if the file was successfully deleted
	if resp.StatusCode == http.StatusOK {
		return c.JSON(http.StatusOK, map[string]bool{"deleted": true})
	}

	return c.String(http.StatusInternalServerError, "Failed to delete file")
}

type FileResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	FileName  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

func retrieveFileHandler(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")
	fileID := c.Param("file_id")

	// Create the GET request
	url := fmt.Sprintf("https://api.openai.com/v1/files/%s", fileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to retrieve the file information
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve file information")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Parse the JSON response
	var fileResponse FileResponse
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse JSON")
	}

	return c.JSON(http.StatusOK, fileResponse)
}

func RetrieveFile(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")
	fileID := c.Param("file_id")

	// Create the GET request
	url := fmt.Sprintf("https://api.openai.com/v1/files/%s/content", fileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create request")
	}

	// Set the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+openAPIKey)

	// Make the request to retrieve the file content
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve file content")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response")
	}

	// Write the file content to a file
	fileName := fmt.Sprintf("%s.jsonl", fileID) // You can change the file extension if needed
	err = ioutil.WriteFile(fileName, body, 0644)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save file content")
	}

	return c.String(http.StatusOK, "File content saved successfully")
}
