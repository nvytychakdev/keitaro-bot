package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/nvytychakdev/keitaro-bot/app"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		slog.Debug("Failed to load environment file, looking for global env...")
	}

	app.Execute()
}
