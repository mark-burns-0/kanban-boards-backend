package api

import (
	"backend/internal/auth"
	"backend/internal/board"
	"backend/internal/platform/config"
	"backend/internal/platform/http"
	"backend/internal/platform/storage/postgres"
	"backend/internal/platform/validation"
	"backend/internal/user"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Storage interface {
	Close() error
}

func Run() error {
	validator := validation.New()
	config := config.NewConfig()
	config.SetLogger()
	storage, err := postgres.NewStorage(config)
	if err != nil {
		slog.Error("Failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// repositories
	userRepo := user.NewUserRepository(storage)
	authRepo := auth.NewAuthRepository(storage)

	//services
	userService := user.NewUserService(userRepo, config)
	authService := auth.NewAuthService(authRepo, config)

	// handlers
	handlers := http.Handlers{
		AuthHandler:  auth.NewAuthHandler(validator, authService),
		UserHandler:  user.NewUserHandler(validator, userService),
		BoardHandler: board.NewBoardHandler(validator),
	}

	app := http.NewApp(config)
	http.RegisterRoutes(app, handlers)
	return runServer(config, app, storage)
}

func runServer(config *config.Config, app *fiber.App, storage Storage) error {
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- app.Listen(
			fmt.Sprintf(":%s", config.GetServerPort()),
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
		_, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()
		if err := app.Shutdown(); err != nil {
			slog.Error(
				"Shutdown error",
				slog.String("error", err.Error()),
			)
		}
		if err := storage.Close(); err != nil {
			slog.Error(
				"closing db",
				slog.String("error", err.Error()),
			)
		}
		return nil
	case err := <-srvErr:
		return err
	}
}
