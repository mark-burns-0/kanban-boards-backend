package api

import (
	"backend/internal/auth"
	"backend/internal/board"
	"backend/internal/card"
	"backend/internal/comment"
	"backend/internal/platform/config"
	"backend/internal/platform/http"
	"backend/internal/platform/lang"
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
	lang := lang.NewRegistry()
	config := config.NewConfig()
	config.SetLogger()
	storage, err := postgres.NewStorage(config)
	if err != nil {
		slog.Error("Failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// repositories
	userRepo, err := user.NewUserRepository(storage)
	if err != nil {
		return err
	}
	authRepo := auth.NewAuthRepository(storage)
	boardRepo := board.NewBoardRepository(storage)
	cardRepo := card.NewCardRepository(storage)
	commentRepo := comment.NewCommentRepository(storage)

	//services
	userService := user.NewUserService(userRepo, config)
	authService := auth.NewAuthService(authRepo, config)
	cardService := card.NewCardService(cardRepo)
	boardService := board.NewBoardService(boardRepo, cardService)
	commentService := comment.NewCommentService(commentRepo)

	// handlers
	handlers := http.Handlers{
		AuthHandler:    auth.NewAuthHandler(validator, lang, authService),
		UserHandler:    user.NewUserHandler(validator, lang, userService),
		BoardHandler:   board.NewBoardHandler(validator, lang, boardService),
		CardHandler:    card.NewCardHandler(validator, lang, cardService),
		CommentHandler: comment.NewCommentHandler(validator, lang, commentService),
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
