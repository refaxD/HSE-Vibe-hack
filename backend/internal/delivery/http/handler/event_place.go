package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
)

type EventPlaceHandler struct {
	usecase domain.EventPlaceUsecase
}

func NewEventPlaceHandler(uc domain.EventPlaceUsecase) *EventPlaceHandler {
	return &EventPlaceHandler{usecase: uc}
}

// GetAll godoc
// @Summary      Get all event places
// @Tags         EventPlaces
// @Produce      json
// @Success      200  {array}  domain.EventPlace
// @Failure      500  {object}  apperrors.AppError
// @Router       /event-places/all [get]
func (h *EventPlaceHandler) GetAll(c *fiber.Ctx) error {
	events, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(events)
}

// GetByID godoc
// @Summary      Get event place by ID
// @Tags         EventPlaces
// @Produce      json
// @Param        event_id  path      string  true  "Event Place ID"
// @Success      200       {object}  domain.EventPlace
// @Failure      400       {object}  apperrors.AppError
// @Failure      404       {object}  apperrors.AppError
// @Router       /event-places/id/{event_id} [get]
func (h *EventPlaceHandler) GetByID(c *fiber.Ctx) error {
	event, err := h.usecase.GetByID(c.Context(), c.Params("event_id"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(event)
}
