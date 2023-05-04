package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken string
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func main() {
	configuration := createConfig()
	fmt.Println(configuration.TelegramBotToken)
	serv, pass, err := ValidatePass("set Linux 123456")
	if err != nil {
		panic(err)
	}
	fmt.Printf("serv: %v, %v\n", serv, pass)

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Text {
			case "open":
				del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID+1)
				msg.ReplyMarkup = numericKeyboard
				go DeleteNextMsg(bot, del)
			}

			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			// del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID+1)
			// go DeleteNextMsg(bot, del)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// switch update.Message.Command() {
		// case "start":
		// 	msg.Text = "Добро пожаловать в моего ботика"
		// case "set":
		// 	pass := update.Message.From.String()
		// 	ValidatePass(pass)
		// 	msg.Text = "Сервис добавлен"
		// case "del":
		// 	msg.Text = "Сервис Удален"
		// case "get":
		// 	msg.Text = "Пароль 123456"
		// default:
		// 	msg.Text = "I don't know that command"
		// }

		// if _, err := bot.Send(msg); err != nil {
		// 	panic(err)
		// }
		// del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID)

		// go DeleteNextMsg(bot, del)

	}
}

func ValidatePass(s string) (service string, pass string, err error) {
	s = strings.TrimSpace(s)
	tmp := strings.Split(s, " ")
	if len(tmp) != 3 {
		err = errors.New("wrong parametrs")
		return
	}
	service = tmp[1]
	pass = tmp[2]
	fmt.Printf("s: %v\n", s)
	return
}

func DeleteNextMsg(bot *tgbotapi.BotAPI, del tgbotapi.DeleteMessageConfig) {
	// err := errors.New("msg does't not exist")
	time.Sleep(10 * time.Second)
	bot.Send(del)

}

func createConfig() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
		// log.Panic(err)
	}
	return configuration
}
