package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"OpenAI-api/api/request"
)

const (
	urlChat         = "https://api.openai.com/v1/chat/completions"
	urlCompletions  = "https://api.openai.com/v1/chat/completions"
	urlEmbeddings   = "https://api.openai.com/v1/embeddings"
	urlImageCreate  = "https://api.openai.com/v1/images/generations"
	urlImageEdit    = "https://api.openai.com/v1/images/edits"
	urlImageVariate = "https://api.openai.com/v1/images/variations"
)

func HandleChat(c echo.Context) error {
	return processChatRequest(c, urlChat)
}

func HandleCompletions(c echo.Context) error {
	return processCompletionsRequest(c, urlCompletions)
}

func HandleEmbeddings(c echo.Context) error {
	return processEmbeddingsRequest(c, urlEmbeddings)
}

func HandleImageCreate(c echo.Context) error {
	return processImageCreate(c, urlImageCreate)
}

func HandleImageEdit(c echo.Context) error {
	return processImageEdit(c, urlImageEdit)
}

func HandleImageVariate(c echo.Context) error {
	return processImageVariation(c, urlImageVariate)
}

func processChatRequest(c echo.Context, url string) error {
	requestBody, err := request.GetChatRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func processCompletionsRequest(c echo.Context, url string) error {
	requestBody, err := request.GetCompletionsRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func processEmbeddingsRequest(c echo.Context, url string) error {
	requestBody, err := request.GetEmbeddingsRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func processImageCreate(c echo.Context, url string) error {
	requestBody, err := request.GetImageCreateRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func processImageEdit(c echo.Context, url string) error {
	requestBody, err := request.GetImageEditRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}

func processImageVariation(c echo.Context, url string) error {
	requestBody, err := request.GetImageVariationRequestBody(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body, err := request.ProcessRequest(url, requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}
