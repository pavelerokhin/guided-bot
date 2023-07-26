package embeddings

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

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
