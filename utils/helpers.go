package utils

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// GetLogger returns the global zerolog logger
func GetLogger() zerolog.Logger {
	return Logger
}

// joinClauses joins a slice of strings with the given separator.
func JoinClauses(clauses []string, sep string) string {
	if len(clauses) == 0 {
		return ""
	}
	out := clauses[0]
	for i := 1; i < len(clauses); i++ {
		out += sep + clauses[i]
	}
	return out
}

// itoa converts an int to string (avoiding strconv import for this file).
func Itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// generateUUID generates a simple UUID-like string
func GenerateUUID() string {
	return uuid.New().String()
}
