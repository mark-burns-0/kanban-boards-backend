package middleware

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func parseAllowedMethodRoutes(routes []fiber.Route) []struct{ path, method string } {
	allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	var parsedSlice []struct{ path, method string }
	for _, route := range routes {
		if slices.Contains(allowedMethods, route.Method) && route.Path != "/" {
			parsedSlice = append(parsedSlice, struct{ path, method string }{
				path:   route.Path,
				method: route.Method,
			})
		}
	}
	return parsedSlice
}

func NotFound(c *fiber.Ctx) error {
	app := c.App()
	routes := parseAllowedMethodRoutes(app.GetRoutes())

	if !contains(routes, c) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": "API endpoint not found",
			"path":   c.Path(),
		})
	}
	return c.Next()
}

func contains(routes []struct{ path, method string }, ctx *fiber.Ctx) bool {
	currentPath := ctx.Path()
	currentMethod := ctx.Method()
	currentParts := strings.Split(strings.Trim(currentPath, "/"), "/")

	for _, route := range routes {
		if route.method != currentMethod {
			continue
		}

		routeParts := strings.Split(strings.Trim(route.path, "/"), "/")

		if len(currentParts) != len(routeParts) {
			continue
		}

		match := true
		for i := range len(routeParts) {
			if strings.HasPrefix(routeParts[i], ":") {
				continue
			}
			if routeParts[i] != currentParts[i] {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}
	return false
}
