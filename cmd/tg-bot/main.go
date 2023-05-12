package main

import (
	"log"
	"tg-Golang-VK/internal/pkg/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}

}
