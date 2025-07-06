package constants

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT                string
	OPENAI_API_KEY      string
	GMAIL_CLIENT_ID     string
	GMAIL_CLIENT_SECRET string
	GMAIL_REDIRECT_URI  string
	DB_URL              string
)

const (
	RESUME_CATEGORY = "resumes"
)

func init() {
	// Load env vars from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system env vars")
	}

	PORT = getEnv("PORT", "8080")
	OPENAI_API_KEY = getEnv("OPENAI_API_KEY", "")
	GMAIL_CLIENT_ID = getEnv("GMAIL_CLIENT_ID", "")
	GMAIL_CLIENT_SECRET = getEnv("GMAIL_CLIENT_SECRET", "")
	GMAIL_REDIRECT_URI = getEnv("GMAIL_REDIRECT_URI", "")
	DB_URL = getEnv("DB_URL", "coldemailer.db")
}

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

func PrintENV() {
	log.Printf("PORT: %s", PORT)
	log.Printf("OPENAI_API_KEY: %s", OPENAI_API_KEY)
	log.Printf("GMAIL_CLIENT_ID: %s", GMAIL_CLIENT_ID)
	log.Printf("GMAIL_CLIENT_SECRET: %s", GMAIL_CLIENT_SECRET)
	log.Printf("GMAIL_REDIRECT_URI: %s", GMAIL_REDIRECT_URI)
	log.Printf("DB_URL: %s", DB_URL)
}
