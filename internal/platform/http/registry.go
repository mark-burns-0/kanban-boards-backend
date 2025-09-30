package http

import "backend/internal/platform/http/routes"

type Handlers struct {
	routes.AuthHandler
	routes.ChallengeHandler
	routes.NotificationHandler
}
