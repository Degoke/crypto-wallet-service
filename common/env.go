package common

import (
	"log"

	"github.com/joho/godotenv"
)

var myENV map[string]string

func LoadENV() {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	myENV = env
}

func GetEnv(key string) string {
	return myENV[key]
}
