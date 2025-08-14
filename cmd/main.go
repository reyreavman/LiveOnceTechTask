package main

import (
	"liveoncetechtask/internal/server"
	"log"
	"os"
)

func main() {
	port := os.Getenv("APP_PORT")

	app := server.NewApp()

	if err := app.Run(port); err != nil {
		log.Fatal("%w", err.Error())
	}
}
