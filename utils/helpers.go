package utils

import "fmt"

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
