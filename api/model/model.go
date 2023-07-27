package model

type RequestBody interface {
	ChatRequestBody | CompletionsRequestBody | EmbeddingsRequestBody
}
