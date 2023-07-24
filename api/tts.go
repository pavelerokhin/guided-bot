package api

import "github.com/labstack/echo/v4"

type TTSRequest struct {
	Model string `json:"model" yaml:"model"`
	Input string `json:"input" yaml:"input"`
}

func TTS(c echo.Context) error {
	return nil
}
