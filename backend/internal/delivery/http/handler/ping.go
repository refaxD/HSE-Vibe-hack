package handler

import "github.com/gofiber/fiber/v2"

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping godoc
// @Summary      Health check
// @Tags         Ping
// @Produce      json
// @Success      200  {string}  string  "success"
// @Router       /ping/ping [get]
func (h *PingHandler) Ping(c *fiber.Ctx) error {
	return c.JSON("success")
}
