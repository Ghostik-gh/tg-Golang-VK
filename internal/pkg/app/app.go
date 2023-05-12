package app

import (
	"log"
	"os"
	"encoding/json"
	"tg-Golang-VK/internal/botapi"
)

type Config struct {
	TelegramBotToken string
}

func Run() error {
	configuration := mustConfig()
	botapi.StartDB()
	bot := botapi.New(configuration.TelegramBotToken)
	bot.StartBot()
	return nil
}

// Считывает конфиг в котором находится токен для бота
func mustConfig() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	log.Println(configuration.TelegramBotToken)
	if err != nil {
		log.Fatal("Can't configurate: create config.json with yours telegram token")
	}
	return configuration
}
