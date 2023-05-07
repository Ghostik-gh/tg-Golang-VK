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
	_ "github.com/lib/pq"
)

func main() {

	configuration := mustConfig()
	fmt.Println(configuration.TelegramBotToken)

	StartDB()
	StartBot(configuration)

}

func Inline(chatID int64, username string) tgbotapi.InlineKeyboardMarkup {

	data, err := UserData(int(chatID), username)
	if len(data) == 0 {
		rows := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Для добавления сервиса используйте set", "Use /set"))
		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows)
		return numericKeyboard
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, v := range data {
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v.service, tgbotapi.EscapeText("MarkdownV2", "`"+v.password+"`")))
		rows = append(rows, row)
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return numericKeyboard
}

func ValidatePass(s string) (service string, pass string, err error) {
	s = strings.TrimSpace(s)
	tmp := strings.Split(s, " ")
	if len(tmp) != 2 {
		err = errors.New("wrong parametrs")
		return
	}
	service = tmp[0]
	pass = tmp[1]
	fmt.Printf("VALIDATE PASS: %v\n", s)
	return
}
func ValidateGet(s string) (service string, err error) {
	s = strings.TrimSpace(s)
	tmp := strings.Split(s, " ")
	if len(tmp) != 1 {
		err = errors.New("wrong parametrs")
		return
	}
	service = tmp[0]
	return
}

func DeleteNextMsg(bot *tgbotapi.BotAPI, chatID int64, MessageID int) {
	del := tgbotapi.NewDeleteMessage(chatID, MessageID)
	err := errors.New("msg does't not exist")
	fmt.Printf("err: %v\n", err)
	time.Sleep(10 * time.Second)
	bot.Send(del)

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
