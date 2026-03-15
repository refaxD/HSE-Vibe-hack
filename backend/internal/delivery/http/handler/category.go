package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
)

type CategoryHandler struct {
	usecase domain.CategoryUsecase
}

func NewCategoryHandler(uc domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: uc}
}

// GetAll godoc
// @Summary      Get all categories
// @Tags         Categories
// @Produce      json
// @Success      200  {array}  domain.Category
// @Failure      500  {object}  apperrors.AppError
// @Router       /categories/all [get]
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	categories, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(categories)
}

// GetByPlaceID godoc
// @Summary      Get categories by place ID
// @Tags         Categories
// @Produce      json
// @Param        place_id  path      string  true  "Place ID"
// @Success      200       {array}   domain.Category
// @Failure      400       {object}  apperrors.AppError
// @Failure      404       {object}  apperrors.AppError
// @Router       /categories/id/place/{place_id} [get]
func (h *CategoryHandler) GetByPlaceID(c *fiber.Ctx) error {
	placeID := c.Params("place_id")
	categories, err := h.usecase.GetByPlaceID(c.Context(), placeID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(categories)
}

func handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
	}
	return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
}
