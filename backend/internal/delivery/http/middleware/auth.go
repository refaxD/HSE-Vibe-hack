package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
)

func AuthMiddleware(sessionRepo domain.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("auth")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer: ")
		if token == authHeader {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		session, err := sessionRepo.Get(c.Context(), token)
		if err != nil || session == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var ctx domain.SessionContext
		if err := json.Unmarshal([]byte(session), &ctx); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		c.Locals("sessionContext", ctx)
		return c.Next()
	}
}
