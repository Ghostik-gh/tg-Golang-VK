package main

import (
	"database/sql"
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
	time.Sleep(3 * time.Second)
	fmt.Printf("os.Getenv(dbInfo): %v\n", os.Getenv(strings.Join(strings.Split(dbInfo, " "), "")))
	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := createTable(); err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(5 * time.Second)

	collectData("Ghostik322", 1, "Hello, World", []string{"answer_1", "answer_2"})
	time.Sleep(5 * time.Second)

	fmt.Println(getNumberOfUsers())
	fmt.Println(getAll())

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

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Command() {
			case "help":
				msg.Text = "/set /get /del"
			case "start":
				msg.ReplyMarkup = numericKeyboard
				// del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID+1)
				// go DeleteNextMsg(bot, del)
			case "get":
				// getAll()
				msg.Text, err = getAll()
				if err != nil {
					panic(err)
				}
			case "del":

			case "set":
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

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

func createTable() error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Создаем таблицу users
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}

	return nil
}

func getNumberOfUsers() (int64, error) {

	var count int64

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//Отправляем запрос в БД для подсчета числа уникальных пользователей
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getAll() (string, error) {

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return "-1", err
	}
	defer db.Close()
	var name, answer string
	row := db.QueryRow("SELECT USERNAME, answer FROM users;")
	err = row.Scan(&name, &answer)
	if err != nil {
		return "-1", err
	}

	return name + ": " + answer, nil
}

// Собираем данные полученные ботом
func collectData(username string, chatid int64, message string, answer []string) error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Конвертируем срез с ответом в строку
	answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}
