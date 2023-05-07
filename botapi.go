package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(config Config) {

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	set := false
	del := false
	for update := range updates {
		// Обработка получения пароля не через аргументы
		if del && update.Message != nil {
			del = false
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			serv, err := ValidateGet(update.Message.Text)
			if err != nil {
				continue
			}
			err = DelPassword(update.Message.From.UserName, msg.ChatID, serv)
			if err == nil {
				msg.Text = "Удалил сервис: " + serv
			}
			continue
		}
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
			go DeleteNextMsg(bot, msg.ChatID, update.Message.MessageID)
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
/del - Удаляет серевис из хранилища
Бот обрабаывает только команды и пару слов после set, все остальное делается через кнопки под сообщениями`
			case "start":
				msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
				msg.Text = `Для того чтобы увидеть пароль кликните на сервис
Справочник по командам
/set - Используется для установления пароля
/get - Выводит список всех сервисов
/del - Удаляет серевис из хранилища
Бот обрабаывает только команды и пару слов после set, все остальное делается через кнопки под сообщениями`
			case "get":
				msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
				msg.Text = "Для того чтобы увидеть пароль кликните на сервис\nПароль удалится через 10 секунд"
			case "del":
				del = true
				msg.Text = "Введите название сервиса"
				// service, err := ValidateGet(update.Message.Text)
				// if err != nil {
				// 	fmt.Printf("err.Error(): %v\n", err)
				// 	continue
				// }
				// err = DelPassword(update.Message.From.UserName, msg.ChatID, service)
				// if err == nil {
				// 	msg.Text = "Удалил сервис: " + service
				// }
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

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

			go DeleteNextMsg(bot, msg.ChatID, update.CallbackQuery.Message.MessageID+1)
			edit := tgbotapi.NewEditMessageReplyMarkup(msg.ChatID, update.CallbackQuery.Message.MessageID,
				tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("/get", "/get"))))
			if _, err := bot.Send(edit); err != nil {
				fmt.Printf("err: %v\n", err)
			}

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
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v.service, tgbotapi.EscapeText("MarkdownV2", "`"+v.password+"`")))
		rows = append(rows, row)
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return numericKeyboard
}

func DeleteNextMsg(bot *tgbotapi.BotAPI, chatID int64, MessageID int) {
	del := tgbotapi.NewDeleteMessage(chatID, MessageID)
	err := errors.New("msg does't not exist")
	fmt.Printf("err: %v\n", err)
	time.Sleep(10 * time.Second)
	bot.Send(del)

}
