package embeddings

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
	url = "https://api.openai.com/v1/embeddings"
)

func CreateEmbeddings(c echo.Context) error {
	return processEmbeddingsRequest(c, url)
}

func processEmbeddingsRequest(c echo.Context, url string) error {
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

func getRequestBody(c echo.Context) (*model.EmbeddingsRequestBody, error) {
	requestBody := model.EmbeddingsRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return nil, err
	}

	if requestBody.Model == "" || requestBody.Input == nil {
		return nil, errors.New("required parameters are not set (required: Model, Input)")
	}

	return &requestBody, nil
}
