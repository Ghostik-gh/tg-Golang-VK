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

type Config struct {
	TelegramBotToken string
}

type RowServ struct {
	service  string
	password string
}

func main() {
	configuration := mustConfig()
	fmt.Println(configuration.TelegramBotToken)

	time.Sleep(5 * time.Second)
	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := createTable(); err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}
	}

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	set := false
	for update := range updates {
		// offset := update.Message.MessageID

		// Обработка получения пароля не через аргументы
		if set && update.Message != nil {
			set = false
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			serv, pass, err := ValidatePass(update.Message.Text)
			if err != nil {
				continue
			}
			err = AddPassword(update.Message.From.UserName, msg.ChatID, serv, pass)
			if err == nil {
				msg.Text = "Добавил пароль от " + serv
			}
			del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID)
			go DeleteNextMsg(bot, del)
			if _, err := bot.Send(msg); err != nil {
				fmt.Printf("err: %v\n", err)
			}
			continue
		}

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Command() {
			case "help":
				msg.Text = `Справочник по командам
/set - Используется для установления пароля
/get - Выводит список всех сервисов
/del - Удаляет серевер из хранилища
Бот обрабаывает только команды и пару слов после set, все остальное делается через кнопки под сообщениями`
			case "start":
				msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
				msg.Text = "Для того чтобы увидеть пароль кликните на сервис\nПароль удалится через 10 секунд"
			case "get":
				msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
				msg.Text = "Для того чтобы увидеть пароль кликните на сервис\nПароль удалится через 10 секунд"
			case "del":
				service, err := ValidateGet(update.Message.Text)
				if err != nil {
					fmt.Printf("err.Error(): %v\n", err)
					continue
				}
				err = DelPassword(update.Message.From.UserName, msg.ChatID, service)
				if err == nil {
					msg.Text = "Удалил сервис: " + service
				}
			case "set":
				set = true
				msg.Text = "Напишите сервис и пароль к нему через пробел"
			default:
				msg.Text = "Я не знаю такой команды"
			}
			if _, err := bot.Send(msg); err != nil {
				fmt.Printf("err: %v\n", err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Пароль удаляется через 10 секунд")
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			del := tgbotapi.NewDeleteMessage(msg.ChatID, update.CallbackQuery.Message.MessageID+1)
			go DeleteNextMsg(bot, del)
			edit := tgbotapi.NewEditMessageReplyMarkup(msg.ChatID, update.CallbackQuery.Message.MessageID,
				tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("/get", "/get"))))
			// Попробовать не паниковать)
			if _, err := bot.Send(edit); err != nil {
				fmt.Printf("err: %v\n", err)
			}
			// tgbotapi.ReplyKeyboardRemove()

		}

	}
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
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v.service, "<code>"+v.password+"</code>"))
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
	if len(tmp) != 2 {
		err = errors.New("wrong parametrs")
		return
	}
	service = tmp[1]
	return
}

func DeleteNextMsg(bot *tgbotapi.BotAPI, del tgbotapi.DeleteMessageConfig) {
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

// Интерфейс для получения одного пароля
// fmt.Printf("update.Message.Text: %v\n", update.Message.Text)
// service, err := ValidateGet(update.Message.Text)
// if err != nil {
// 	fmt.Printf("err.Error(): %v\n", err)
// 	continue
// }
// pass, err := Password(int(msg.ChatID), service, update.Message.From.UserName)
// msg.Text = "Ваш пароль " + pass + "\nЧepeз 10 секуд пароль будет удален"
// del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID+1)
// go DeleteNextMsg(bot, del)
// if err != nil {
// 	fmt.Printf("err.Error(): %v\n", err)
// 	continue
// }
