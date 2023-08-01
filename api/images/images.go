package images

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"OpenAI-api/api/model"
	"OpenAI-api/api/request"
)

const (
	urlCreate  = "https://api.openai.com/v1/images/generations"
	urlEdit    = "https://api.openai.com/v1/images/edits"
	urlVariate = "https://api.openai.com/v1/images/variations"
)

func HandleCreate(c echo.Context) error {
	return processImageCreate(c, urlCreate)
}

func HandleEdit(c echo.Context) error {
	return processImageEdit(c, urlEdit)
}

func HandleVariate(c echo.Context) error {
	return processImageVariate(c, urlVariate)
}

func processImageCreate(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	requestBody, err := getCreateRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := request.MakeRequest(requestBody, url, apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getCreateRequestBody(c echo.Context) (*model.ImageCreateRequestBody, error) {
	requestBody := model.ImageCreateRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Prompt == "" {
		return nil, errors.New("required parameters are not set (required: Prompt)")
	}

	return &requestBody, nil
}

func processImageEdit(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	requestBody, err := getEditRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := request.MakeRequest(requestBody, url, apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getEditRequestBody(c echo.Context) (*model.ImageEditRequestBody, error) {
	requestBody := model.ImageEditRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Image == "" || requestBody.Prompt == "" {
		return nil, errors.New("required parameters are not set (required: Image, Prompt)")
	}

	return &requestBody, nil
}

func processImageVariate(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	requestBody, err := getVariationRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req, err := request.MakeRequest(requestBody, url, apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.SendRequest(nil, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func getVariationRequestBody(c echo.Context) (*model.ImageVariateRequestBody, error) {
	requestBody := model.ImageVariateRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Image == "" {
		return nil, errors.New("required parameters are not set (required: Image)")
	}

	return &requestBody, nil
}
