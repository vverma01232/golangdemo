package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvMongoURI fetches Environment Info for connecting to Mongo DB
func LoadEnv() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file due to: ", err.Error())

	}
}
