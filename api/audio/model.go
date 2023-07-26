package audio

type RequestBody struct {
	File           string  `json:"file"`
	Model          string  `json:"model"`  // ID of the model to use. Only whisper-1 is currently available.
	Prompt         string  `json:"prompt"` // If set should match audio language
	ResponseFormat string  `json:"response_format"`
	Temperature    float32 `json:"temperature"`
	Language       string  `json:"language"` // The language of the input audio. Supplying the input language in ISO-639-1 format will improve accuracy and latency.
}
