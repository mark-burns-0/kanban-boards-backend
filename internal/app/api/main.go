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
	"backend/internal/infrastructure/config"
	"backend/internal/infrastructure/http"
	"backend/internal/infrastructure/lang"
	"backend/internal/infrastructure/storage/postgres"
	"backend/internal/infrastructure/validation"
	userDomain "backend/internal/user/domain"
	userRepository "backend/internal/user/repository"
	userTransport "backend/internal/user/transport"
	"database/sql"
	"sync"
	"time"

	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

const (
	shutdownTimeout = 10 * time.Second
)

type Storage interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	GetDB() *sql.DB
	Close() error
}

type Config interface {
	GetAppName() string
	GetServerPort() string
	GetBcryptPower() string
	GetAccessTokenTTL() string
	GetRefreshTokenTTL() string
	GetAccessTokenSecret() string
	GetRefreshTokenSecret() string
}

type App struct {
	validator *validation.Validator
	lang      *lang.Registry
	app       *fiber.App
	storage   Storage
	config    Config
}

func NewApp() (*App, error) {
	validator := validation.New()
	lang := lang.NewRegistry()
	cfg := config.NewConfig()
	cfg.SetLogger()
	storage, err := postgres.NewStorage(cfg)
	if err != nil {
		slog.Error("Failed to create storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("create storage: %w", err)
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
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := a.app.Shutdown(); err != nil {
			shutdownErr <- fmt.Errorf("http server shutdown: %w", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.storage.Close(); err != nil {
			shutdownErr <- fmt.Errorf("storage close: %w", err)
			return
		}
	}()

	go func() {
		defer close(shutdownErr)
		wg.Wait()
	}()

	var errs []error
	for {
		select {
		case err, ok := <-shutdownErr:
			if !ok {
				if len(errs) > 0 {
					return fmt.Errorf("shutdown completed with errors: %v", errs)
				}
				slog.Info("Graceful shutdown completed")
				return nil
			}
			if err != nil {
				errs = append(errs, err)
				slog.Error("Shutdown error", slog.Any("error", err))
			}
		case <-ctx.Done():
			return fmt.Errorf("shutdown timeout: %w", ctx.Err())
		}
	}
}
