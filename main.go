package main

import (
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Infolog = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)

func main() {
	configuration := mustConfig()
	StartDB()
	StartBot(configuration)
}

func mustConfig() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Can't configurate: create config.json with yours telegram token")
	}
	return configuration
}
