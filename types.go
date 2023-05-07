package main

type Config struct {
	TelegramBotToken string
}

// Хранит пару из сервиса и пароля к нему
type RowServ struct {
	service  string
	password string
}
