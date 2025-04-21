package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotenvValue(desiredValue string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	value := os.Getenv(desiredValue)
	return value
}
