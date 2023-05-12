package botapi

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

// Хранит пару из сервиса и пароля к нему
type RowServ struct {
	service  string
	password string
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

// Стартует работу базы данных
func StartDB() {
	// Задержка пока разворачивается Postgres
	time.Sleep(6 * time.Second)
	if os.Getenv("CREATE_TABLE") == "yes" {
		if os.Getenv("DB_SWITCH") == "on" {
			if err := createTable(); err != nil {
				Infolog.Printf("err: %v\n", err)
			}
		}
	}
}

// Подключается к БД
// Желательно в каждой функции закрывать соединение через close()
func connectDB() (db *sql.DB) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		Infolog.Printf("DATABASE ERROR: %v\n", err)
	}
	return
}

// Созадет таблицу если она еще не создана
func createTable() error {
	db := connectDB()
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE passwords(id SERIAL PRIMARY KEY, username TEXT, chat_id INT, service TEXT, password TEXT);`); err != nil {
		Infolog.Printf("createTable: %v\n", err)
		return err
	}
	return nil
}

// Добавляет одну запись
func AddPassword(username string, chatid int64, service string, password string) error {
	db := connectDB()
	defer db.Close()
	data := `INSERT INTO passwords(username, chat_id, service, password) VALUES($1, $2, $3, $4);`
	if _, err := db.Exec(data, username, chatid, service, Encrypt(password)); err != nil {
		Infolog.Printf("AddPassword: %v\n", err)
		return err
	}
	return nil
}

// Удаляет одну запись из таблицы
func DelPassword(username string, chatid int64, service string) error {
	db := connectDB()
	defer db.Close()
	data := `DELETE FROM passwords WHERE service = $1;`
	if _, err := db.Exec(data, service); err != nil {
		Infolog.Printf("DelPassword: %v\n", err)
		return err
	}
	return nil
}

// Получает все пары из сервисов и паролей
func UserData(chat_id int, username string) ([]RowServ, error) {
	db := connectDB()
	defer db.Close()
	var service_name, pass string
	var data []RowServ
	rows, err := db.Query(`SELECT service, password FROM passwords WHERE chat_id = $1 AND username = $2;`, chat_id, username)
	if err != nil {
		Infolog.Printf("UserData: %v\n", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&service_name, &pass); err != nil {
			return data, err
		}

		data = append(data, RowServ{service_name, Decrypt(pass)})
	}

	if err = rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}
