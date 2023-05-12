package app

import (
	"fmt"
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
	// file, _ := os.Open("config.json")
	// decoder := json.NewDecoder(file)
	configuration := Config{}
	// configuration.TelegramBotToken = os.Getenv("TOKEN")
	configuration.TelegramBotToken = "6201240293:AAHdMaz4Qu8ShdCfCQSFskqWGe1-bxw1-uU"
	fmt.Printf("configuration.TelegramBotToken: %v\n", configuration.TelegramBotToken)
	// err := decoder.Decode(&configuration)
	// if err != nil {
	// 	log.Fatal("Can't configurate: create config.json with yours telegram token")
	// }
	return configuration
}
