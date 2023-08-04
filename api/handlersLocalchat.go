package api

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"

	"OpenAI-api/api/model"
	"OpenAI-api/api/request"
)

func HandleChatL(prompt string) (string, error) {
	body, err := processChatRequestL(prompt, urlChat)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HandleCompletionsL(prompt string) (string, error) {
	body, err := processCompletionsRequestL(prompt, urlCompletions)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HandleEmbeddingsL(prompt string) (string, error) {
	body, err := processEmbeddingsRequestL(prompt, urlEmbeddings)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HandleImageCreateL(prompt string) ([]string, error) {
	body, err := processImageCreateL(prompt, urlImageCreate)
	if err != nil {
		return nil, err
	}

	imageResponse := model.ImageResponse{}
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, image := range imageResponse.Data {
		out = append(out, image.B64Json)
	}

	return out, nil
}

func HandleImageEditL(image, prompt string) ([]string, error) {
	body, err := processImageEditL(image, prompt, urlImageEdit)
	if err != nil {
		return nil, err
	}

	imageResponse := model.ImageResponse{}
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, image := range imageResponse.Data {
		out = append(out, image.Url)
	}

	return out, nil
}

func HandleImageVariateL(image string) ([]string, error) {
	body, err := processImageVariationL(image, urlImageVariate)
	if err != nil {
		return nil, err
	}

	imageResponse := model.ImageResponse{}
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, image := range imageResponse.Data {
		out = append(out, image.Url)
	}

	return out, nil
}

func processChatRequestL(prompt string, url string) ([]byte, error) {
	requestBody := model.ChatRequestBody{
		Model: viper.GetString("openAI.model"),
		// TODO: FINISH OTHER REQUESTED FIELDS
	}
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processCompletionsRequestL(prompt string, url string) ([]byte, error) {
	requestBody := model.CompletionsRequestBody{
		Model:  viper.GetString("openAI.model"),
		Prompt: prompt,
	}

	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processEmbeddingsRequestL(prompt, url string) ([]byte, error) {
	requestBody := model.CompletionsRequestBody{
		Model:  viper.GetString("openAI.model"),
		Prompt: prompt,
	}

	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processImageCreateL(prompt, url string) ([]byte, error) {
	requestBody := model.ImageCreateRequestBody{
		Prompt:         prompt,
		Size:           viper.GetString("openAI.imageSize"),
		N:              viper.GetInt64("openAI.imageN"),
		ResponseFormat: viper.GetString("openAI.imageResponseFormat"),
	}

	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processImageEditL(imageFile, prompt, url string) ([]byte, error) {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	// open image file
	image, err := OpenFileAndReadAllLines(imageFile)
	if err != nil {
		return nil, err
	}

	requestBody := model.ImageEditRequestBody{
		Image:  image,
		Prompt: prompt,
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func processImageVariationL(imageFile, url string) ([]byte, error) {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return nil, errors.New("OpenAI API key not found")
	}

	// open image file
	image, err := OpenFileAndReadAllLines(imageFile)
	if err != nil {
		return nil, err
	}

	requestBody := model.ImageVariateRequestBody{
		Image: image,
	}

	req, err := request.MakeRequest(&requestBody, url, apiKey)
	if err != nil {
		return nil, err
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func OpenFileAndReadAllLines(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Join all the lines into a single string
	result := strings.Join(lines, "\n")

	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}
