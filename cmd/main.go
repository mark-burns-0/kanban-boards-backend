package main

import (
	"backend/internal/app/api"
	"log"
	"log/slog"
	"os"
)

func main() {
	const op = "cmd.main"
	app, err := api.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		slog.Error("app shutdown", slog.String("op", op), slog.Any("err", err))
		os.Exit(1)
	}
}
