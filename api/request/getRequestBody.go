package request

import (
	"OpenAI-api/api/model"
	"errors"
	"github.com/labstack/echo/v4"
)

func GetChatRequestBody(c echo.Context) (*model.ChatRequestBody, error) {
	requestBody := model.ChatRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Messages == nil {
		return nil, errors.New("required parameters are not set (required: Model, Messages)")
	}

	return &requestBody, nil
}

func GetCompletionsRequestBody(c echo.Context) (*model.CompletionsRequestBody, error) {
	requestBody := model.CompletionsRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Prompt == nil {
		return nil, errors.New("required parameters are not set (required: Model, Prompt)")
	}

	return &requestBody, nil
}

func GetImageCreateRequestBody(c echo.Context) (*model.ImageCreateRequestBody, error) {
	requestBody := model.ImageCreateRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Prompt == "" {
		return nil, errors.New("required parameters are not set (required: Prompt)")
	}

	return &requestBody, nil
}

func GetImageEditRequestBody(c echo.Context) (*model.ImageEditRequestBody, error) {
	requestBody := model.ImageEditRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Image == "" || requestBody.Prompt == "" {
		return nil, errors.New("required parameters are not set (required: Image, Prompt)")
	}

	return &requestBody, nil
}

func GetImageVariationRequestBody(c echo.Context) (*model.ImageVariateRequestBody, error) {
	requestBody := model.ImageVariateRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Image == "" {
		return nil, errors.New("required parameters are not set (required: Image)")
	}

	return &requestBody, nil
}

func GetEmbeddingsRequestBody(c echo.Context) (*model.CompletionsRequestBody, error) {
	requestBody := model.CompletionsRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Prompt == nil {
		return nil, errors.New("required parameters are not set (required: Model, Prompt)")
	}

	return &requestBody, nil
}
