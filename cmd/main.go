package main

import (
	"liveoncetechtask/internal/server"
	"log"
)

func main() {
	port := "8080"

	app := server.NewApp()

	if err := app.Run(port); err != nil {
		log.Fatal("%w", err.Error())
	}
}
