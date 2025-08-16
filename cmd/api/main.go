package main

import (
	"os"

	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/Agero19/AnnotateX-api/internal/db"
	"github.com/Agero19/AnnotateX-api/internal/logger"
	"github.com/Agero19/AnnotateX-api/internal/repository"
	"github.com/Agero19/AnnotateX-api/internal/server"
	"github.com/joho/godotenv"
)

func main() {

	//load env by godotenv
	_ = godotenv.Load()
	// This is the entry point for the API server.
	// The server will start listening on the address specified in the .env file.
	// It will also connect to the database using the connection string provided in the .env file.
	// Additional configurations such as maximum open connections, idle connections, and idle time can be set as well.

	// Load the configuration
	cfg := config.LoadConfig()

	// Initialize the logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting AnnotateX API", "env", cfg.Env, "port", cfg.Port)
	log.Debug("Debug level logging enabled")

	// Initialize the database connection
	db, err := db.New(
		cfg.DB.URL,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)

	if err != nil {
		log.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repository
	repo := repository.NewRepository(db)

	// var user_sample repository.User
	// user_sample.Username = "testuser"
	// user_sample.Email = "testuser@example.com"
	// user_sample.Password = "securepassword"

	// // Create a sample user
	// if err := repo.Users.Create(&user_sample); err != nil {
	// 	log.Error("Failed to create user", "error", err)
	// }
	// log.Info("Created user", "id", user_sample.ID)

	// Start the API server
	app := server.NewApp(cfg, repo, log)
	mux := app.Mount()
	//TODO: Implement graceful shutdown for the server
	if err := app.Run(mux); err != nil {
		log.Error("Failed to run API server", "error", err)
	}
}
