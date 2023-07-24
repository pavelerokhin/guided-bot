package api

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

type ChatCompletionRequest struct {
	Model        string             `json:"model"`
	Messages     []Message          `json:"messages"`
	Functions    []Function         `json:"functions,omitempty"`
	Stop         interface{}        `json:"stop,omitempty"`
	MaxTokens    int                `json:"max_tokens,omitempty"`
	PresencePen  float64            `json:"presence_penalty,omitempty"`
	FrequencyPen float64            `json:"frequency_penalty,omitempty"`
	LogitBias    map[string]float64 `json:"logit_bias,omitempty"`
	User         string             `json:"user,omitempty"`
	Temperature  float64            `json:"temperature,omitempty"`
	TopP         float64            `json:"top_p,omitempty"`
	N            int                `json:"n,omitempty"`
	Stream       bool               `json:"stream,omitempty"`
	FunctionCall string             `json:"function_call,omitempty"`
}

type Message struct {
	Role     string       `json:"role"`
	Content  string       `json:"content"`
	Name     string       `json:"name,omitempty"`
	Function FunctionCall `json:"function_call,omitempty"`
}

type Function struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Parameters  FunctionParam `json:"parameters"`
}

type FunctionParam struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type FunctionCall struct {
	Name string      `json:"name"`
	Args interface{} `json:"args,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// https://platform.openai.com/docs/api-reference/completions
func Completion(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var completionReq ChatCompletionRequest
	err := c.Bind(&completionReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(completionReq)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var completionResp ChatCompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, completionResp)
}

func CompletionStreaming(c echo.Context) error {
	return nil
}
