package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	configuration := mustConfig()
	fmt.Println(configuration.TelegramBotToken)

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
