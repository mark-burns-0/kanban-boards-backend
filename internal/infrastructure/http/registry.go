package http

import "backend/internal/infrastructure/http/routes"

type Handlers struct {
	routes.AuthHandler
	routes.BoardHandler
	routes.UserHandler
	routes.CardHandler
	routes.CommentHandler
}
