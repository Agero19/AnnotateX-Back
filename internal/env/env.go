package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return fallback
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if v, err := time.ParseDuration(val); err == nil {
		return v
	}
	return fallback
}
