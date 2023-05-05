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

var (
	host     = os.Getenv("HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname   = os.Getenv("DBNAME")
	sslmode  = os.Getenv("SSLMODE")
)

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

// var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
// 		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
// 		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
// 		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
// 		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
// 	),
// )

func main() {
	configuration := createConfig()
	fmt.Println(configuration.TelegramBotToken)

	time.Sleep(4 * time.Second)
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
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			var err error
			switch update.Message.Command() {
			case "help":
				msg.Text = "/set /get /del"
			// TODO: Inline и Keyboard
			// case "start":
			// 	msg.ReplyMarkup = numericKeyboard
			case "get":
				fmt.Printf("update.Message.Text: %v\n", update.Message.Text)
				service, err := ValidateGet(update.Message.Text)
				if err != nil {
					fmt.Printf("err.Error(): %v\n", err)
					continue
				}
				pass, err := GetPassword(int(msg.ChatID), service, update.Message.From.UserName)
				msg.Text = "Ваш пароль " + pass + `
Через 10 секуд пароль будет удален`
				del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID+1)
				go DeleteNextMsg(bot, del)
				if err != nil {
					fmt.Printf("err.Error(): %v\n", err)
					continue
				}
			case "del":
				service, err := ValidateGet(update.Message.Text)
				if err != nil {
					fmt.Printf("err.Error(): %v\n", err)
					continue
				}
				DelPassword(update.Message.From.UserName, msg.ChatID, service)
				msg.Text = "Удалил сервис: " + service
			case "set":
				serv, pass, err := ValidatePass(update.Message.Text)
				if err != nil {
					continue
				}
				AddPassword(update.Message.From.UserName, msg.ChatID, serv, pass)
				msg.Text = "Добавил пароль от " + serv
				del := tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID)
				go DeleteNextMsg(bot, del)
			default:
				msg.Text = "я не знаю такой команды"
			}
			if _, err = bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

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
		log.Panic(err)
	}
	return configuration
}

func createTable() error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Создаем таблицу users
	if _, err = db.Exec(`CREATE TABLE passwords(id SERIAL PRIMARY KEY, username TEXT, chat_id INT, service TEXT, password TEXT);`); err != nil {
		return err
	}

	return nil
}

func getNumberOfUsers() (int64, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var count int64
	// Количество записей в таблице
	row := db.QueryRow("SELECT COUNT(DISTINCT id) FROM passwords;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetPassword(chat_id int, service, username string) (string, error) {

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return "Can't Connect", err
	}
	defer db.Close()
	var service_name, pass string
	row := db.QueryRow(`SELECT service, password FROM passwords WHERE service = $1;`, service) // WHERE chat_id = $1 AND service = $2 AND username = $3;`, chat_id, service, username)
	err = row.Scan(&service_name, &pass)
	if err != nil {
		return "Not Found", err
	}
	return service_name + ": " + pass, nil
}

// Собираем данные полученные ботом
func AddPassword(username string, chatid int64, service string, password string) error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Конвертируем срез с ответом в строку
	// answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO passwords(username, chat_id, service, password) VALUES($1, $2, $3, $4);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, `@`+username, chatid, service, password); err != nil {
		return err
	}

	return nil
}

func DelPassword(username string, chatid int64, service string) error {
	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	data := `DELETE FROM passwords WHERE service = $1;`
	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, service); err != nil {
		return err
	}
	return nil
}
