package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	dir, err := os.Getwd()

	if err != nil {
		log.Println("Error getting current directory")
		panic(err)
	}
	dir = dir + "/.env"

	log.Println(" --> Loading .env file from: ", dir)

	err = godotenv.Load(dir)
	if err != nil {
		log.Println(" --> Error loading .env file")
	}
}
