package model

// request

type ChatRequestBody struct {
	Model            string             `json:"model"`
	Messages         []Message          `json:"messages"`
	Functions        []Function         `json:"functions,omitempty"`
	FunctionCall     string             `json:"function_call,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`       // default 1
	TopP             float64            `json:"top_p,omitempty"`             // default 1
	N                int64              `json:"n,omitempty"`                 // default 1
	Stream           bool               `json:"stream,omitempty"`            // default false
	Stop             string             `json:"stop,omitempty"`              // string or array, default nil
	MaxTokens        int64              `json:"max_tokens,omitempty"`        // default infinity
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`  // from -2.0 to 2.0
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"` // from -2.0 to 2.0
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}

type Message struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	Name         string `json:"name,omitempty"`
	FunctionCall string `json:"function_call,omitempty"`
}

type Function struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Parameters   string `json:"parameters"`
	FunctionCall string `json:"function_call"`
}

// response

type ChatResponse struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int64      `json:"created"`
	Choices []Choice   `json:"choices"`
	Usage   TokenUsage `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
