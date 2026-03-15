package http

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/hse-vibe-hack/backend/internal/delivery/http/handler"
	"github.com/hse-vibe-hack/backend/internal/delivery/http/middleware"
	"github.com/hse-vibe-hack/backend/internal/domain"
)

func SetupRoutes(
	app *fiber.App,
	pingH *handler.PingHandler,
	userH *handler.UserHandler,
	categoryH *handler.CategoryHandler,
	placeH *handler.PlaceHandler,
	eventH *handler.EventPlaceHandler,
	pathH *handler.PathHandler,
	sessionRepo domain.SessionRepository,
) {
	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Ping
	ping := app.Group("/ping")
	ping.Get("/ping", pingH.Ping)

	// Auth
	auth := app.Group("/auth")
	auth.Post("/register", userH.Register)
	auth.Post("/login", userH.Login)
	auth.Get("/profile", middleware.AuthMiddleware(sessionRepo), userH.Profile)

	// Categories
	categories := app.Group("/categories")
	categories.Get("/all", categoryH.GetAll)
	categories.Get("/id/place/:place_id", categoryH.GetByPlaceID)

	// Places
	places := app.Group("/places")
	places.Get("/all", placeH.GetAll)
	places.Get("/id/:place_id", placeH.GetByID)
	places.Get("/id/category/:category_id", placeH.GetByCategoryID)

	// Event Places
	eventPlaces := app.Group("/event-places")
	eventPlaces.Get("/all", eventH.GetAll)
	eventPlaces.Get("/id/:event_id", eventH.GetByID)

	// Paths
	path := app.Group("/path")
	path.Get("/all", pathH.GetAll)
	path.Get("/:path_id", pathH.GetByID)
	path.Get("/:path_id/segments", pathH.GetSegments)
	path.Post("/create", pathH.Create)
}
