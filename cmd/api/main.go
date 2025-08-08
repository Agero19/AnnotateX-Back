package main

import (
	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/Agero19/AnnotateX-api/internal/logger"
	"github.com/Agero19/AnnotateX-api/internal/server"
)

func main() {
	// This is the entry point for the API server.
	// The server will start listening on the address specified in the .env file.
	// It will also connect to the database using the connection string provided in the .env file.
	// Additional configurations such as maximum open connections, idle connections, and idle time can be set as well.

	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting AnnotateX API", "env", cfg.Env, "port", cfg.Port)
	log.Debug("Debug level logging enabled")

	// Start the API server
	app := server.NewApp(cfg)
	mux := app.Mount()
	//TODO: Implement graceful shutdown for the server
	if err := app.Run(mux); err != nil {
		log.Error("Failed to run API server", "error", err)
	}
}
