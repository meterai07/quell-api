package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVar() {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load .env file!")
	}
}
