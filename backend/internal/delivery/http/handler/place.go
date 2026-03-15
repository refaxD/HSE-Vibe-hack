package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
)

type PlaceHandler struct {
	usecase domain.PlaceUsecase
}

func NewPlaceHandler(uc domain.PlaceUsecase) *PlaceHandler {
	return &PlaceHandler{usecase: uc}
}

// GetAll godoc
// @Summary      Get all places
// @Tags         Places
// @Produce      json
// @Success      200  {array}  domain.Place
// @Failure      500  {object}  apperrors.AppError
// @Router       /places/all [get]
func (h *PlaceHandler) GetAll(c *fiber.Ctx) error {
	places, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(places)
}

// GetByID godoc
// @Summary      Get place by ID
// @Tags         Places
// @Produce      json
// @Param        place_id  path      string  true  "Place ID"
// @Success      200       {object}  domain.Place
// @Failure      400       {object}  apperrors.AppError
// @Failure      404       {object}  apperrors.AppError
// @Router       /places/id/{place_id} [get]
func (h *PlaceHandler) GetByID(c *fiber.Ctx) error {
	place, err := h.usecase.GetByID(c.Context(), c.Params("place_id"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(place)
}

// GetByCategoryID godoc
// @Summary      Get places by category ID
// @Tags         Places
// @Produce      json
// @Param        category_id  path      string  true  "Category ID"
// @Success      200          {array}   domain.Place
// @Failure      400          {object}  apperrors.AppError
// @Failure      404          {object}  apperrors.AppError
// @Router       /places/id/category/{category_id} [get]
func (h *PlaceHandler) GetByCategoryID(c *fiber.Ctx) error {
	places, err := h.usecase.GetByCategoryID(c.Context(), c.Params("category_id"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(places)
}
