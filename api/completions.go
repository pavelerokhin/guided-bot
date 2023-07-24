package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type CompletionLegacyRequest struct {
	Model            string             `json:"model"`
	Prompt           interface{}        `json:"prompt"`
	Suffix           string             `json:"suffix,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Logprobs         int                `json:"logprobs,omitempty"`
	Echo             bool               `json:"echo,omitempty"`
	Stop             interface{}        `json:"stop,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	BestOf           int                `json:"best_of,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}

type CompletionLegacyResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   TokenUsage         `json:"usage"`
}

type CompletionChoice struct {
	Text         string         `json:"text"`
	Index        int            `json:"index"`
	Logprobs     []TokenLogprob `json:"logprobs"`
	FinishReason string         `json:"finish_reason"`
}

type TokenLogprob struct {
	TokenIndex int     `json:"token_index"`
	Token      string  `json:"token"`
	Logprob    float64 `json:"logprob"`
}

// https://platform.openai.com/docs/api-reference/completions
func CreateCompletion(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var completionReq CompletionLegacyRequest
	err := c.Bind(&completionReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(completionReq)
	if err != nil {
		return err
	}

	url := "https://api.openai.com/v1/completions"

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var completionResp CompletionLegacyResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, completionResp)
}
