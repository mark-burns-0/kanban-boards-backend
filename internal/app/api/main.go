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
	"backend/internal/shared/ports/repository"
	userDomain "backend/internal/user/domain"
	userRepository "backend/internal/user/repository"
	userTransport "backend/internal/user/transport"
	"time"

	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

const (
	shutdownTimeout = 10 * time.Second
)

type Storage interface {
	Close() error
}

type App struct {
	config    *config.Config
	validator *validation.Validator
	lang      *lang.Registry
	app       *fiber.App
	storage   repository.Storage
}

func NewApp() (*App, error) {
	validator := validation.New()
	lang := lang.NewRegistry()
	cfg := config.NewConfig()
	cfg.SetLogger()
	storage, err := postgres.NewStorage(cfg)
	if err != nil {
		slog.Error("Failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return &App{
		config:    cfg,
		validator: validator,
		lang:      lang,
		storage:   storage,
	}, nil
}

func (a *App) setupServices() (http.Handlers, error) {
	// repositories
	userRepo, err := userRepository.NewUserRepository(a.storage)
	if err != nil {
		return http.Handlers{}, err
	}
	authRepo := userRepository.NewAuthRepository(a.storage)
	boardRepo := boardRepository.NewBoardRepository(a.storage)
	cardRepo := cardRepository.NewCardRepository(a.storage)
	commentRepo := commentRepository.NewCommentRepository(a.storage)

	//services
	userService := userDomain.NewUserService(userRepo, a.config)
	authService := userDomain.NewAuthService(authRepo, a.config)
	cardService := cardDomain.NewCardService(cardRepo, boardRepo)
	boardService := boardDomain.NewBoardService(boardRepo, cardService)
	commentService := commentDomain.NewCommentService(commentRepo)

	// handlers
	return http.Handlers{
		AuthHandler:    userTransport.NewAuthHandler(a.validator, a.lang, authService),
		UserHandler:    userTransport.NewUserHandler(a.validator, a.lang, userService),
		BoardHandler:   boardTransport.NewBoardHandler(a.validator, a.lang, boardService),
		CardHandler:    cardTransport.NewCardHandler(a.validator, a.lang, cardService),
		CommentHandler: commentTransport.NewCommentHandler(a.validator, a.lang, commentService),
	}, nil
}

func (a *App) Run() error {
	handlers, err := a.setupServices()
	if err != nil {
		return err
	}
	a.app = http.NewApp(a.config)
	http.RegisterRoutes(a.app, handlers)

	return a.runServer()
}

func (a *App) runServer() error {
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- a.app.Listen(
			fmt.Sprintf(":%s", a.config.GetServerPort()),
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
		return a.gracefulShutdown()
	case err := <-srvErr:
		return fmt.Errorf("server error: %w", err)
	}
}

func (a *App) gracefulShutdown() error {
	slog.Info("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	shutdownErr := make(chan error, 2)

	go func() {
		if err := a.app.Shutdown(); err != nil {
			shutdownErr <- fmt.Errorf("http server shutdown: %w", err)
			return
		}
		shutdownErr <- nil
	}()

	go func() {
		if err := a.storage.Close(); err != nil {
			shutdownErr <- fmt.Errorf("storage close: %w", err)
			return
		}
		shutdownErr <- nil
	}()

	for i := 0; i < 2; i++ {
		select {
		case err := <-shutdownErr:
			if err != nil {
				slog.Error("Shutdown error", "error", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("shutdown timeout: %w", ctx.Err())
		}
	}
	slog.Info("Graceful shutdown completed")
	return nil
}
