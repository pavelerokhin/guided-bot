package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type ChatCompletionRequest struct {
	Model            string             `json:"model"`
	Messages         []ChatMessage      `json:"messages"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Stop             interface{}        `json:"stop,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Choices []ChatChoice `json:"choices"`
	Usage   TokenUsage   `json:"usage"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func CreateChatCompletion(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var chatReq ChatCompletionRequest
	err := c.Bind(&chatReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return err
	}

	url := "https://api.openai.com/v1/chat/completions"

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

	var chatResp ChatCompletionResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, chatResp)
}
