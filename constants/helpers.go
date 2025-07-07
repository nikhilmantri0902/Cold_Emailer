package constants

import (
	"fmt"
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		if fallback == "" {
			log.Printf("WARNING: Env var %s not set and no fallback given", key)
		}
		return fallback
	}
	return value
}

func getEnvFloat32(key string, fallback float32) float32 {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return fallback
	}
	var f float32
	_, err := fmt.Sscanf(value, "%f", &f)
	if err != nil {
		return fallback
	}
	return f
}

func getEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return fallback
	}
	var i int
	_, err := fmt.Sscanf(value, "%d", &i)
	if err != nil {
		return fallback
	}
	return i
}
