package constants

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// environment variables
var (
	PORT                         string
	OPENAI_API_KEY               string
	GMAIL_CLIENT_ID              string
	GMAIL_CLIENT_SECRET          string
	GMAIL_REDIRECT_URI           string
	OPENAI_MODEL                 string
	OPENAI_TEMPERATURE           float32
	OPENAI_MAX_COMPLETION_TOKENS int
)

// constants
const (
	RESUME_CATEGORY = "resumes"
)

var (
	PG_DB_URL = fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", "coldemailer", "coldemailer", "coldemailer")
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
	OPENAI_MODEL = getEnv("OPENAI_MODEL", "gpt-3.5-turbo")
	OPENAI_TEMPERATURE = getEnvFloat32("OPENAI_TEMPERATURE", 0.7)
	OPENAI_MAX_COMPLETION_TOKENS = getEnvInt("OPENAI_MAX_COMPLETION_TOKENS", 512)
}

func PrintENV() {
	log.Printf("PORT: %s", PORT)
	log.Printf("OPENAI_API_KEY: %s", OPENAI_API_KEY)
	log.Printf("GMAIL_CLIENT_ID: %s", GMAIL_CLIENT_ID)
	log.Printf("GMAIL_CLIENT_SECRET: %s", GMAIL_CLIENT_SECRET)
	log.Printf("GMAIL_REDIRECT_URI: %s", GMAIL_REDIRECT_URI)
}
