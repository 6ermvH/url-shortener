package main

import (
	"log"

	"github.com/6ermvH/url-shortener/cmd/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
