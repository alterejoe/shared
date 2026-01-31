package env

import (
	"fmt"
	"os"
	"strings"
)

// generic helper to get an environment variable by prefix
func getEnv(prefix, key string) string {
	return os.Getenv(strings.ToUpper(fmt.Sprintf("AUTH0_%s_%s", prefix, key)))
}
func cleanScopes(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	// allow comma or space separation
	split := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	var scopes []string
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s != "" {
			scopes = append(scopes, s)
		}
	}
	return scopes
}

func RequiredEnvs(required []string, prefix string) error {
	for _, key := range required {
		if getEnv(prefix, key) == "" {
			return fmt.Errorf("missing env var AUTH0_%s_%s", prefix, key)
		}
	}
	return nil
}
