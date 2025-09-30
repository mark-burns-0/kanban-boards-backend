package api

import (
	"backend/internal/challenge"
	"backend/internal/notification"
	"backend/internal/platform/http"
	"backend/internal/user"
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	srvErr := make(chan error, 1)

	app := http.NewApp()
	handlers := http.Handlers{
		AuthHandler:         user.NewAuthHandler(),
		ChallengeHandler:    challenge.NewChallengeHandler(),
		NotificationHandler: notification.NewNotificationHandler(),
	}
	http.RegisterRoutes(app, handlers)
	go func() { srvErr <- app.Listen(":8080") }()

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
