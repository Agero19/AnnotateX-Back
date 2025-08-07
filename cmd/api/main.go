package main

import (
	"log"

	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/Agero19/AnnotateX-api/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// This is the entry point for the API server.
	// The server will start listening on the address specified in the .env file.
	// It will also connect to the database using the connection string provided in the .env file.
	// Additional configurations such as maximum open connections, idle connections, and idle time can be set as well.

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
	}

	// Start the API server
	app := server.NewApp(cfg)
	mux := app.Mount()
	//TODO: Implement graceful shutdown for the server
	if err := app.Run(mux); err != nil {
		log.Fatal(err)
	}
}
