package main

import (
	"errors"
	"fmt"
	"strings"
)

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
