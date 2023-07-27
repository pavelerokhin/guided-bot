package complitions

import (
	"OpenAI-api/api/model"
	"OpenAI-api/api/request"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

const (
	url = "https://api.openai.com/v1/completions"
)

func CreateCompletion(c echo.Context) error {
	return processChatRequest(c, url)
}

func processChatRequest(c echo.Context, url string) error {
	apiKey := viper.GetString("openAI.apiKey")
	if apiKey == "" && strings.HasPrefix(url, "https://api.openai.com/") {
		return echo.NewHTTPError(http.StatusUnauthorized, "OpenAI API key not found")
	}

	requestBody, err := getRequestBody(c)
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

func getRequestBody(c echo.Context) (*model.CompletionsRequestBody, error) {
	requestBody := model.CompletionsRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Prompt == nil {
		return nil, errors.New("required parameters are not set (required: Model, Prompt)")
	}

	return &requestBody, nil
}
