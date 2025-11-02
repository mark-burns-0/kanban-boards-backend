package api

import (
	boardDomain "backend/internal/board/domain"
	boardRepository "backend/internal/board/repository"
	boardTransport "backend/internal/board/transport"
	cardDomain "backend/internal/card/domain"
	cardRepository "backend/internal/card/repository"
	cardTransport "backend/internal/card/transport"
	commentDomain "backend/internal/comment/domain"
	commentRepository "backend/internal/comment/repository"
	commentTransport "backend/internal/comment/transport"
	"backend/internal/platform/config"
	"backend/internal/platform/http"
	"backend/internal/platform/lang"
	"backend/internal/platform/storage/postgres"
	"backend/internal/platform/validation"
	userDomain "backend/internal/user/domain"
	userRepository "backend/internal/user/repository"
	userTransport "backend/internal/user/transport"
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
	userRepo, err := userRepository.NewUserRepository(storage)
	if err != nil {
		return err
	}
	authRepo := userRepository.NewAuthRepository(storage)
	boardRepo := boardRepository.NewBoardRepository(storage)
	cardRepo := cardRepository.NewCardRepository(storage)
	commentRepo := commentRepository.NewCommentRepository(storage)

	//services
	userService := userDomain.NewUserService(userRepo, config)
	authService := userDomain.NewAuthService(authRepo, config)
	cardService := cardDomain.NewCardService(cardRepo, boardRepo)
	boardService := boardDomain.NewBoardService(boardRepo, cardService)
	commentService := commentDomain.NewCommentService(commentRepo)

	// handlers
	handlers := http.Handlers{
		AuthHandler:    userTransport.NewAuthHandler(validator, lang, authService),
		UserHandler:    userTransport.NewUserHandler(validator, lang, userService),
		BoardHandler:   boardTransport.NewBoardHandler(validator, lang, boardService),
		CardHandler:    cardTransport.NewCardHandler(validator, lang, cardService),
		CommentHandler: commentTransport.NewCommentHandler(validator, lang, commentService),
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
