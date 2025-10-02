package api

import (
	"backend/internal/challenge"
	"backend/internal/notification"
	"backend/internal/platform/config"
	"backend/internal/platform/http"
	"backend/internal/platform/lang"
	"backend/internal/user"
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
)

func Run() error {
	srvErr := make(chan error, 1)

	app := http.NewApp()
	langRegistry := lang.NewRegistry()
	validator := validator.New()

	config := config.NewConfig()
	config.SetLogger()

	handlers := http.Handlers{
		AuthHandler: user.NewAuthHandler(
			validator,
			langRegistry,
		),
		ChallengeHandler:    challenge.NewChallengeHandler(),
		NotificationHandler: notification.NewNotificationHandler(),
	}
	http.RegisterRoutes(app, handlers)
	go func() { srvErr <- app.Listen(fmt.Sprintf(":%s", config.ServerPort)) }()

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
