package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/hse-vibe-hack/backend/docs"
	httpDelivery "github.com/hse-vibe-hack/backend/internal/delivery/http"
	"github.com/hse-vibe-hack/backend/internal/delivery/http/handler"
	mongoRepo "github.com/hse-vibe-hack/backend/internal/repository/mongo"
	redisRepo "github.com/hse-vibe-hack/backend/internal/repository/redis"
	"github.com/hse-vibe-hack/backend/internal/usecase"
	"github.com/hse-vibe-hack/backend/pkg/config"
)

// @title           HSE Vibe Hack API
// @version         1.0
// @description     Backend API for route generation service
// @host            localhost:3000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name auth
// @description     Enter your bearer token: "Bearer: <token>"

func main() {
	// Load config
	loadEnvFile()
	cfg := config.Load()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}
	log.Println("connected to MongoDB")

	db := mongoClient.Database("apidb")

	// Connect to Redis
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("connected to Redis")

	// Repositories
	userRepo := mongoRepo.NewUserRepository(db)
	categoryRepo := mongoRepo.NewCategoryRepository(db)
	placeRepo := mongoRepo.NewPlaceRepository(db)
	eventPlaceRepo := mongoRepo.NewEventPlaceRepository(db)
	pathRepo := mongoRepo.NewPathRepository(db)
	pathSegmentRepo := mongoRepo.NewPathSegmentRepository(db)
	sessionRepo := redisRepo.NewSessionRepository(redisClient)

	// Usecases
	userUC := usecase.NewUserUsecase(userRepo, sessionRepo, cfg.RedisTTL)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo, placeRepo)
	placeUC := usecase.NewPlaceUsecase(placeRepo, categoryRepo)
	eventPlaceUC := usecase.NewEventPlaceUsecase(eventPlaceRepo)
	pathUC := usecase.NewPathUsecase(pathRepo, pathSegmentRepo, placeRepo, eventPlaceRepo, categoryRepo, cfg)

	// Handlers
	pingH := handler.NewPingHandler()
	userH := handler.NewUserHandler(userUC)
	categoryH := handler.NewCategoryHandler(categoryUC)
	placeH := handler.NewPlaceHandler(placeUC)
	eventH := handler.NewEventPlaceHandler(eventPlaceUC)
	pathH := handler.NewPathHandler(pathUC)

	// Parse entities if configured
	if cfg.ParseEntities {
		parseService := usecase.NewParseService(categoryRepo, placeRepo, eventPlaceRepo, cfg)
		go parseService.Run(context.Background())
	}

	// Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowMethods:     "GET,POST,PUT,OPTIONS",
		AllowCredentials: true,
	}))

	// Routes
	httpDelivery.SetupRoutes(app, pingH, userH, categoryH, placeH, eventH, pathH, sessionRepo)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Printf("server started on port: http://localhost:%s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	app.Shutdown()
}

func loadEnvFile() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	for _, line := range splitLines(string(data)) {
		if line == "" || line[0] == '#' {
			continue
		}
		for i := 0; i < len(line); i++ {
			if line[i] == '=' {
				key := line[:i]
				value := line[i+1:]
				if os.Getenv(key) == "" {
					os.Setenv(key, value)
				}
				break
			}
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
