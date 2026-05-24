package main

import (
	"log"

	"awesomeProject/internal/app"
	"awesomeProject/internal/config"
)

func main() {
	cfg := config.Load()
	application := app.New(cfg)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
