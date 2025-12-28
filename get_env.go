package ptti

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(errors.New("no env found"))
	}

	result := os.Getenv(key)

	if result == "" {
		log.Fatalf("%s does not exist", key)
	}

	return result
}
