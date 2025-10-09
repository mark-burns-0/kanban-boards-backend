package api

import (
	"backend/internal/auth"
	"backend/internal/board"
	"backend/internal/platform/config"
	"backend/internal/platform/http"
	"backend/internal/platform/validation"
	"backend/internal/user"
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	srvErr := make(chan error, 1)

	app := http.NewApp()
	validator := validation.New()
	config := config.NewConfig()
	config.SetLogger()

	handlers := http.Handlers{
		AuthHandler:  auth.NewAuthHandler(validator),
		BoardHandler: board.NewBoardHandler(validator),
		UserHandler:  user.NewUserHandler(validator),
	}

	http.RegisterRoutes(app, handlers)

	go func() {
		srvErr <- app.Listen(
			fmt.Sprintf(":%s", config.ServerPort),
		)
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		// shutdownCtx
		_, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()
		if err := app.Shutdown(); err != nil {
			slog.Error("Shutdown error", err)
		}
		return nil
	case err := <-srvErr:
		return err
	}
}
