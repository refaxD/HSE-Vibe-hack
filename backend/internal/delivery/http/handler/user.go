package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(uc domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

// Register godoc
// @Summary      Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.UserRegisterDTO  true  "Registration data"
// @Success      200   {string}  string                  "Bearer token"
// @Failure      400   {object}  apperrors.AppError
// @Failure      409   {object}  apperrors.AppError
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var dto domain.UserRegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	token, err := h.usecase.Register(c.Context(), dto)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(token)
}

// Login godoc
// @Summary      Login user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.UserLoginDTO  true  "Login credentials"
// @Success      200   {string}  string               "Bearer token"
// @Failure      409   {object}  apperrors.AppError
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var dto domain.UserLoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	token, err := h.usecase.Login(c.Context(), dto)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(token)
}

// Profile godoc
// @Summary      Get current user profile
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  domain.UserPublic
// @Failure      401   {object}  apperrors.AppError
// @Router       /auth/profile [get]
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	ctx, ok := c.Locals("sessionContext").(domain.SessionContext)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	profile, err := h.usecase.Profile(c.Context(), ctx.Email)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(profile)
}
