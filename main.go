package main

import (
	"OpenAI-api/api/audio"
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

	// audio
	e.POST("/audio/transcription", audio.CreateTranscription)
	e.POST("/audio/translation", audio.CreateTranslation)

	if ImageDir != "" {
		e.Static("/assets/images", ImageDir)
	}

	if AudioDir != "" {
		e.Static("/assets/audio", AudioDir)
	}

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
