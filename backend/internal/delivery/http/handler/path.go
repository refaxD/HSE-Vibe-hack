package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hse-vibe-hack/backend/internal/domain"
)

type PathHandler struct {
	usecase domain.PathUsecase
}

func NewPathHandler(uc domain.PathUsecase) *PathHandler {
	return &PathHandler{usecase: uc}
}

// GetAll godoc
// @Summary      Get all paths
// @Tags         Paths
// @Produce      json
// @Success      200  {array}  domain.Path
// @Failure      500  {object}  apperrors.AppError
// @Router       /path/all [get]
func (h *PathHandler) GetAll(c *fiber.Ctx) error {
	paths, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(paths)
}

// GetByID godoc
// @Summary      Get path by ID
// @Tags         Paths
// @Produce      json
// @Param        path_id  path      string  true  "Path ID"
// @Success      200      {object}  domain.Path
// @Failure      400      {object}  apperrors.AppError
// @Failure      404      {object}  apperrors.AppError
// @Router       /path/{path_id} [get]
func (h *PathHandler) GetByID(c *fiber.Ctx) error {
	path, err := h.usecase.GetByID(c.Context(), c.Params("path_id"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(path)
}

// GetSegments godoc
// @Summary      Get path segments by path ID
// @Tags         Paths
// @Produce      json
// @Param        path_id  path      string  true  "Path ID"
// @Success      200      {array}   domain.PathSegment
// @Failure      400      {object}  apperrors.AppError
// @Failure      404      {object}  apperrors.AppError
// @Router       /path/{path_id}/segments [get]
func (h *PathHandler) GetSegments(c *fiber.Ctx) error {
	segments, err := h.usecase.GetSegments(c.Context(), c.Params("path_id"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(segments)
}

// Create godoc
// @Summary      Create a new path (AI-generated route)
// @Tags         Paths
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreatePathDTO  true  "Prompt and start position"
// @Success      200   {array}   domain.Place
// @Failure      400   {object}  apperrors.AppError
// @Failure      500   {object}  apperrors.AppError
// @Router       /path/create [post]
func (h *PathHandler) Create(c *fiber.Ctx) error {
	var dto domain.CreatePathDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.usecase.CreatePath(c.Context(), dto)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(result)
}
