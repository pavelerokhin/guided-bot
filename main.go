package main

import (
	"net/http"
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
	e.POST("/api/endpoint1", endpoint1Handler) // e.g. curl -X POST http://localhost:8080/api/endpoint1
	e.POST("/api/endpoint2", endpoint2Handler)

	e.GET("/version", func(c echo.Context) error {
		return c.String(200, "v0.0.1")
	})

	e.POST("/models/apply", func(c echo.Context) error { return nil })
	e.GET("/models/available", func(c echo.Context) error { return nil })
	e.GET("/models/jobs/:uuid", func(c echo.Context) error { return nil })

	// openAI compatible API endpoint

	// chat
	e.POST("/v1/chat/completions", func(c echo.Context) error { return nil })
	e.POST("/chat/completions", func(c echo.Context) error { return nil })

	// edit
	e.POST("/v1/edits", func(c echo.Context) error { return nil })
	e.POST("/edits", func(c echo.Context) error { return nil })

	// completion
	e.POST("/v1/completions", func(c echo.Context) error { return nil })
	e.POST("/completions", func(c echo.Context) error { return nil })
	e.POST("/v1/engines/:model/completions", func(c echo.Context) error { return nil })

	// embeddings
	e.POST("/v1/embeddings", func(c echo.Context) error { return nil })
	e.POST("/embeddings", func(c echo.Context) error { return nil })
	e.POST("/v1/engines/:model/embeddings", func(c echo.Context) error { return nil })

	// audio
	e.POST("/v1/audio/transcriptions", func(c echo.Context) error { return nil })
	e.POST("/tts", func(c echo.Context) error { return nil })

	// images
	e.POST("/v1/images/generations", func(c echo.Context) error { return nil })

	if ImageDir != "" {
		e.Static("/generated-images", ImageDir)
	}

	if AudioDir != "" {
		e.Static("/generated-audio", AudioDir)
	}

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler for "/api/endpoint1"
func endpoint1Handler(c echo.Context) error {
	// Your logic for endpoint 1 here
	return c.String(http.StatusOK, "Endpoint 1: Success")
}

// Handler for "/api/endpoint2"
func endpoint2Handler(c echo.Context) error {
	// Your logic for endpoint 2 here
	return c.String(http.StatusOK, "Endpoint 2: Success")
}
