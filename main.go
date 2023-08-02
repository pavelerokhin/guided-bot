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
	e.Use(middleware.Static("static"))

	// Routes
	// chat
	e.POST("/v1/chat/completions", api.HandleChat)
	e.POST("/chat/completions", api.HandleChat)

	// completions
	e.POST("/v1/completions", api.HandleCompletions)
	e.POST("/completions", api.HandleCompletions)

	// embeddings
	e.POST("/v1/embeddings", api.HandleEmbeddings)
	e.POST("/embeddings", api.HandleEmbeddings)

	// images
	e.POST("/v1/images/generations", api.HandleImageCreate)
	e.POST("/images/generations", api.HandleImageCreate)
	e.POST("/v1/images/edits", api.HandleImageEdit)
	e.POST("/images/edits", api.HandleImageEdit)
	e.POST("/v1/images/variations", api.HandleImageVariate)
	e.POST("/images/variations", api.HandleImageVariate)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
