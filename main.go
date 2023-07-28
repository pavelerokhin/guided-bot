package main

import (
	"OpenAI-api/api/chat"
	"OpenAI-api/api/completions"
	"OpenAI-api/api/embeddings"
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
	e.POST("/chat/completions", chat.Handle)

	// completions
	e.POST("/v1/completions", completions.Handle)
	e.POST("/completions", completions.Handle)

	// embeddings
	e.POST("/v1/embeddings", embeddings.Handle)
	e.POST("/embeddings", embeddings.Handle)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
