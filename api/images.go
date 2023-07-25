package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type ImageGenerationRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageGenerationResponse struct {
	Created int64         `json:"created"`
	Data    []ImageResult `json:"data"`
}

type ImageResult struct {
	URL string `json:"url"`
}

func CreateImage(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var imageGenReq ImageGenerationRequest
	err := c.Bind(&imageGenReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(imageGenReq)
	if err != nil {
		return err
	}

	url := "https://api.openai.com/v1/images/generations"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var imageGenResp ImageGenerationResponse
	err = json.Unmarshal(body, &imageGenResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, imageGenResp)
}

func CreateImageEdit(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	prompt := c.FormValue("prompt")
	n := 1
	if val := c.FormValue("n"); val != "" {
		n = atoi(val)
	}
	size := "1024x1024"
	if val := c.FormValue("size"); val != "" {
		size = val
	}

	// Process image and mask files from the request
	imageFile, err := c.FormFile("image")
	if err != nil {
		return err
	}

	maskFile, err := c.FormFile("mask")
	if err != nil && err != http.ErrMissingFile {
		return err
	}

	// Prepare request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file
	imagePart, err := writer.CreateFormFile("image", imageFile.Filename)
	if err != nil {
		return err
	}
	imageFileSrc, err := imageFile.Open()
	if err != nil {
		return err
	}
	defer imageFileSrc.Close()
	io.Copy(imagePart, imageFileSrc)

	// Add mask file if provided
	if maskFile != nil {
		maskPart, err := writer.CreateFormFile("mask", maskFile.Filename)
		if err != nil {
			return err
		}
		maskFileSrc, err := maskFile.Open()
		if err != nil {
			return err
		}
		defer maskFileSrc.Close()
		io.Copy(maskPart, maskFileSrc)
	}

	// Add other fields to the request body
	_ = writer.WriteField("prompt", prompt)
	_ = writer.WriteField("n", fmt.Sprintf("%d", n))
	_ = writer.WriteField("size", size)
	_ = writer.WriteField("response_format", "url")
	_ = writer.WriteField("user", "") // Add your end-user identifier here if needed

	// Close the writer and set appropriate headers
	writer.Close()

	url := "https://api.openai.com/v1/images/edits"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response as JSON
	var imageGenResp ImageGenerationResponse
	err = json.Unmarshal(responseBody, &imageGenResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, imageGenResp)
}

func CreateImageVariation(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	n := 1
	if val := c.FormValue("n"); val != "" {
		n = atoi(val)
	}
	size := "1024x1024"
	if val := c.FormValue("size"); val != "" {
		size = val
	}

	// Process image file from the request
	imageFile, err := c.FormFile("image")
	if err != nil {
		return err
	}

	// Prepare request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file
	imagePart, err := writer.CreateFormFile("image", imageFile.Filename)
	if err != nil {
		return err
	}
	imageFileSrc, err := imageFile.Open()
	if err != nil {
		return err
	}
	defer imageFileSrc.Close()
	io.Copy(imagePart, imageFileSrc)

	// Add other fields to the request body
	_ = writer.WriteField("n", fmt.Sprintf("%d", n))
	_ = writer.WriteField("size", size)
	_ = writer.WriteField("response_format", "url")
	_ = writer.WriteField("user", "") // Add your end-user identifier here if needed

	// Close the writer and set appropriate headers
	writer.Close()

	url := "https://api.openai.com/v1/images/variations"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response as JSON
	var imageGenResp ImageGenerationResponse
	err = json.Unmarshal(responseBody, &imageGenResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, imageGenResp)
}
