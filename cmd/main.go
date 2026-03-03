package main

import (
	"log"

	"github.com/6ermvH/url-shortener/cmd/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
