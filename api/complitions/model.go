package complitions

import "OpenAI-api/api/chat"

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
	Usage   chat.TokenUsage    `json:"usage"`
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
