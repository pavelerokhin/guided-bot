package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type EmbeddingsRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"`
	User  string      `json:"user,omitempty"`
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// https://platform.openai.com/docs/api-reference/embeddings
func CreateEmbeddings(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var embeddingsReq EmbeddingsRequest
	err := c.Bind(&embeddingsReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(embeddingsReq)
	if err != nil {
		return err
	}

	url := "https://api.openai.com/v1/embeddings"

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

	var embeddingsResp EmbeddingsResponse
	err = json.Unmarshal(body, &embeddingsResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, embeddingsResp)
}
