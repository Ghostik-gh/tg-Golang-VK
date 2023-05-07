package main

import (
	"errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	helpText = `Справочник по командам
/set - Используется для установления пароля
/get - Выводит список всех сервисов
/del - Удаляет серевис из хранилища
/author - Если у вас есть вопросы или предложения по работе бота пишите мне
Бот обрабаывает только команды и пару слов после set, все остальное делается через кнопки под сообщениями`
)

func StartBot(config Config) {

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	Infolog.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	set := false
	del := false
	get := false
	count := 0
	lastMsg := 0
	for update := range updates {
		// Обработка получения пароля не через аргументы
		if del && update.Message != nil {
			del = false
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			serv, err := ValidateGet(update.Message.Text)
			if err != nil {
				msg.Text = "Не удалось удалить сервис: проверьте правильность написания сервиса используйте /get для просмотра"
				if _, err := bot.Send(msg); err != nil {
					Infolog.Printf("err: %v\n", err)
				}
				continue
			}
			if err := DelPassword(update.Message.From.UserName, msg.ChatID, serv); err == nil {
				msg.Text = "Удалил сервис: " + serv + "\nЕсли пароль не удалился проверьте написание сервиса через /get"
			}
			msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
			if _, err := bot.Send(msg); err != nil {
				Infolog.Printf("err: %v\n", err)
			}
			count = 3
			continue
		}

		// Обработка получения нового сервиса и пароля
		if set && update.Message != nil {
			set = false
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			serv, pass, err := ValidatePass(update.Message.Text)
			if err != nil {
				msg.Text = "Не удалось добавить пароль, данные надо писать в одном сообщении через пробел"
				if _, err := bot.Send(msg); err != nil {
					Infolog.Printf("err: %v\n", err)
				}
				continue
			}
			if err := AddPassword(update.Message.From.UserName, msg.ChatID, serv, pass); err == nil {
				msg.Text = "Добавил пароль от " + serv
			}
			go DeleteMsg(bot, msg.ChatID, update.Message.MessageID)
			// msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
			if _, err := bot.Send(msg); err != nil {
				Infolog.Printf("err: %v\n", err)
			}
			count = 3
			continue
		}

		// Обрабатываем сообщения пользователя
		if update.Message != nil {
			if update.Message.MessageID >= lastMsg {
				lastMsg = update.Message.MessageID
				count = 1
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Command() {
			case "start":
				msg.Text = helpText
			case "help":
				msg.Text = helpText
			case "get":
				count = 0
				get = true
				msg.ReplyMarkup = Inline(update.Message.Chat.ID, update.Message.From.UserName)
				msg.Text = "Для того чтобы увидеть пароль кликните на сервис\nПароль удалится через 10 секунд"
			case "del":
				del = true
				msg.Text = "Введите название сервиса"
			case "set":
				set = true
				msg.Text = "Напишите сервис и пароль к нему через пробел"
			case "author":
				msg.Text = "https://t.me/GhostikGH"
			default:
				msg.Text = "Я не знаю такой команды\n" + helpText
			}
			// Отправляем сообщение
			if _, err := bot.Send(msg); err != nil {
				Infolog.Printf("err: %v\n", err)
			}
		} else if update.CallbackQuery != nil {
			// Обрабатываем нажатия на кнопки под сообщениями
			lastMsg = Max(update.CallbackQuery.Message.MessageID, lastMsg)
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Пароль удаляется через 10 секунд")
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			count++
			if get {
				go DeleteMsg(bot, msg.ChatID, lastMsg+count)
			}
		}
	}
}

// Создает Инлайн клавиатуру под сообщением
func Inline(chatID int64, username string) tgbotapi.InlineKeyboardMarkup {
	data, err := UserData(int(chatID), username)
	if len(data) == 0 {
		rows := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Для добавления сервиса используйте set", "Use /set"))
		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows)
		return numericKeyboard
	}

	if err != nil {
		Infolog.Printf("err: %v\n", err)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	curLen := len(data)
	for i := 0; i < curLen; {
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(data[i].service, data[i].password))
		if curLen-i == 1 {
			row = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(data[i].service, data[i].password))
			i += 1
		} else if curLen-i == 2 {
			row = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(data[i].service, data[i].password),
				tgbotapi.NewInlineKeyboardButtonData(data[i+1].service, data[i+1].password))
			i += 2
		} else if curLen >= 3 {
			row = tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(data[i].service, data[i].password),
				tgbotapi.NewInlineKeyboardButtonData(data[i+1].service, data[i+1].password),
				tgbotapi.NewInlineKeyboardButtonData(data[i+2].service, data[i+2].password))
			i += 3
		}
		rows = append(rows, row)
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return numericKeyboard
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Удаляет переданное сообщение через 5 секунд
func DeleteMsg(bot *tgbotapi.BotAPI, chatID int64, MessageID int) {
	del := tgbotapi.NewDeleteMessage(chatID, MessageID)
	err := errors.New("msg does't not exist")
	Infolog.Printf("err: %v\n", err)
	time.Sleep(5 * time.Second)
	bot.Send(del)
}
