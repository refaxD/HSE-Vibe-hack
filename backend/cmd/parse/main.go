package main

import (
	"context"
	"fmt"
	"log"
	"os"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoRepo "github.com/hse-vibe-hack/backend/internal/repository/mongo"
	"github.com/hse-vibe-hack/backend/internal/usecase"
	"github.com/hse-vibe-hack/backend/pkg/config"
)

func main() {
	loadEnvFile()
	cfg := config.Load()

	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}
	log.Println("connected to MongoDB")

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("connected to Redis")

	db := mongoClient.Database("apidb")
	categoryRepo := mongoRepo.NewCategoryRepository(db)
	placeRepo := mongoRepo.NewPlaceRepository(db)
	eventPlaceRepo := mongoRepo.NewEventPlaceRepository(db)

	parseService := usecase.NewParseService(categoryRepo, placeRepo, eventPlaceRepo, cfg)
	parseService.Run(ctx)

	log.Println("parse finished")
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
