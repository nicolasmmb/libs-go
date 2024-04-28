package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	dir, err := os.Getwd()

	if err != nil {
		log.Println("Erro ao pegar o diretório atual")
		panic(err)
	}
	dir = dir + "/.env"

	log.Println("Buscando arquivo .env em: ", dir)

	err = godotenv.Load(dir)
	if err != nil {
		log.Println("Error loading .env file")
		// panic(err)
	}
}
