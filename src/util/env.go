package util

import (
	"log"
	"os"
	"strconv"
)

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env: %q is not set", key)
	}
	return value
}

func MustGetEnvInt(key string) int {
	value, err := strconv.Atoi(MustGetEnv(key))
	if err != nil {
		log.Fatalf("env: %q is not valid 'int'", key)
	}
	return value
}
