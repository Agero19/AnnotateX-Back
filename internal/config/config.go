package config

import (
	"time"

	"github.com/Agero19/AnnotateX-api/internal/env"
)

type dbConfig struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type Config struct {
	Env  string
	Port string
	DB   dbConfig
	// Another configurations structs if needed
	// cache, logging, s3, auth
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		Env:  env.GetString("ENV", "local"),
		Port: env.GetString("PORT", ":8080"),
		DB: dbConfig{
			URL:          env.GetString("DB_URL", "postgres://user:password@localhost:5432/dbname"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
		},
	}
	// panic if config is not set including fallbacks
	if cfg.Env == "" || cfg.Port == "" || cfg.DB.URL == "" {
		panic("Missing required config values")
	}

	return cfg
}
