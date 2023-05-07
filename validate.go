package main

import (
	"errors"
	"strings"
)

// Проверяет передвнную пару сервис пароль
// Пароль должен содеражть аски символы от 33 до 122
// Количество аргуметов строго равно 2
func ValidatePass(s string) (service string, pass string, err error) {
	s = strings.TrimSpace(s)
	tmp := strings.Split(s, " ")
	if len(tmp) != 2 {
		err = errors.New("wrong parametrs")
		return
	}
	service = tmp[0]
	pass = tmp[1]
	for _, v := range []byte(pass) {
		if int(v) < 33 || int(v) > 122 {
			err = errors.New("unallowed symbols")
			return
		}
	}
	return
}

// Проверяет строку после команды /del
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
