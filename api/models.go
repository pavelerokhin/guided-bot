package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Model struct {
	ID         string   `json:"id"`
	Object     string   `json:"object"`
	OwnedBy    string   `json:"owned_by"`
	Permission []string `json:"permission"`
}

type ModelsResponse struct {
	Data   []Model `json:"data"`
	Object string  `json:"object"`
}

func List(c echo.Context) error {
	c.Logger().Debugf("Listing models from OpenAI API")

	openAPIKey := viper.GetString("openAI.apiKey")
	url := "https://api.openai.com/v1/models"

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
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var modelsResponse ModelsResponse
	err = json.Unmarshal(body, &modelsResponse)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, modelsResponse)
}

func Retrieve(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")
	modelID := c.Param("modelID")
	url := fmt.Sprintf("https://api.openai.com/v1/models/%s", modelID)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var model Model
	err = json.Unmarshal(body, &model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model)
}
