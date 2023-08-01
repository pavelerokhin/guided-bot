package model

type ImageCreateRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int64  `json:"n,omitempty"`               // default 1
	Size           string `json:"string,omitempty"`          // default 1024x1024
	ResponseFormat string `json:"response_format,omitempty"` // default "url"
	User           string `json:"user,omitempty"`
}

type ImageEditRequestBody struct {
	Image          string `json:"image"`
	Prompt         string `json:"prompt"`
	Mask           int64  `json:"mask,omitempty"`
	N              int64  `json:"n,omitempty"`               // default 1
	Size           string `json:"string,omitempty"`          // default 1024x1024
	ResponseFormat string `json:"response_format,omitempty"` // default "url"
	User           string `json:"user,omitempty"`
}

type ImageVariateRequestBody struct {
	Image          string `json:"image"`
	N              int64  `json:"n,omitempty"`               // default 1
	Size           string `json:"string,omitempty"`          // default 1024x1024
	ResponseFormat string `json:"response_format,omitempty"` // default "url"
	User           string `json:"user,omitempty"`
}
