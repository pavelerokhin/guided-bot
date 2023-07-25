package main

import (
	"OpenAI-api/api"
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var (
	ImageDir string
	AudioDir string
)

func main() {
	// Initialize viper and read configurations
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	// Create an Echo instance
	e := echo.New()

	// Set the logger to use a custom format
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetHeader("${time_rfc3339} ${level}")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure()) // https://echo.labstack.com/middleware/secure/ TODO: see if we need it
	e.Use(middleware.Static("static"))

	// Routes
	e.GET("/version", func(c echo.Context) error {
		return c.String(200, "v0.0.1")
	})

	e.GET("/models/list", api.List)
	e.GET("/models/retireve", api.Retrieve)

	// openAI compatible API endpoint
	// chat
	//e.POST("/chat", api.CreateChat)

	// edit
	//e.POST("/edit", api.Edit)

	// completion
	e.POST("/completion", api.CreateCompletion)

	// embeddings
	e.POST("/embedding", api.CreateEmbeddings)

	// audio
	//e.POST("/audio/transcription", api.Transcription)
	//e.POST("/audio/tts", api.TTS)

	// images
	e.POST("/images/generation", api.CreateImage)

	if ImageDir != "" {
		e.Static("/generated-images", ImageDir)
	}

	if AudioDir != "" {
		e.Static("/generated-audio", AudioDir)
	}

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
