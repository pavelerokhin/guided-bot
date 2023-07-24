package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func CreateTranscription(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Save the uploaded file to a temporary location
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create("tempfile.mp3") // Save the file with the appropriate extension (e.g., mp3, wav, etc.)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	model := c.FormValue("model")
	prompt := c.FormValue("prompt")
	responseFormat := c.FormValue("response_format")
	temperature := c.FormValue("temperature")
	language := c.FormValue("language")

	url := "https://api.openai.com/v1/audio/transcriptions"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	// Create a form data payload with the audio file and other parameters
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Add the audio file to the request
	audioFile, err := os.Open("tempfile.mp3") // Use the same file path as the one used for saving the uploaded file
	if err != nil {
		return err
	}
	defer audioFile.Close()

	audioPart, err := writer.CreateFormFile("file", "audio.mp3")
	if err != nil {
		return err
	}

	_, err = io.Copy(audioPart, audioFile)
	if err != nil {
		return err
	}

	// Add other parameters to the request
	writer.WriteField("model", model)
	writer.WriteField("prompt", prompt)
	writer.WriteField("response_format", responseFormat)
	writer.WriteField("temperature", temperature)
	writer.WriteField("language", language)

	err = writer.Close()
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Delete the temporary audio file after processing the request
	os.Remove("tempfile.mp3")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transcription": string(body),
	})
}

func CreateTranslation(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Save the uploaded file to a temporary location
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create("tempfile.m4a") // Save the file with the appropriate extension (e.g., mp3, wav, etc.)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	model := c.FormValue("model")
	prompt := c.FormValue("prompt")
	responseFormat := c.FormValue("response_format")
	temperature := c.FormValue("temperature")

	url := "https://api.openai.com/v1/audio/translations"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	// Create a form data payload with the audio file and other parameters
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Add the audio file to the request
	audioFile, err := os.Open("tempfile.m4a") // Use the same file path as the one used for saving the uploaded file
	if err != nil {
		return err
	}
	defer audioFile.Close()

	audioPart, err := writer.CreateFormFile("file", "audio.m4a")
	if err != nil {
		return err
	}

	_, err = io.Copy(audioPart, audioFile)
	if err != nil {
		return err
	}

	// Add other parameters to the request
	writer.WriteField("model", model)
	writer.WriteField("prompt", prompt)
	writer.WriteField("response_format", responseFormat)
	writer.WriteField("temperature", temperature)

	err = writer.Close()
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Delete the temporary audio file after processing the request
	os.Remove("tempfile.m4a")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"translation": string(body),
	})
}
