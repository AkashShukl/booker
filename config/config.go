package config

import (
	"log"

	"github.com/joho/godotenv"
)

// room number with name,
// number goes in db
var Rooms = map[string]string{
	"1": "Alpha",
	"2": "Beta",
	"3": "Gamma",
}

const ()

var env map[string]string

func LoadEnv() {
	var err error
	env, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
}

func GetEnv() map[string]string {
	if env == nil {
		LoadEnv()
	}
	return env
}
