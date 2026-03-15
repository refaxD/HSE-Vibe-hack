package config

import (
	"os"
	"strconv"
)

type Config struct {
	NodeEnv  string
	Port     string
	MongoURI string

	RedisHost string
	RedisPort string
	RedisTTL  int

	TimepadKey   string
	RouterKey    string
	RouterURL    string
	OpenAIAPIKey string

}

func Load() *Config {
	redisTTL, _ := strconv.Atoi(getEnv("REDIS_TTL", "86400"))

	return &Config{
		NodeEnv:  getEnv("NODE_ENV", "development"),
		Port:     getEnv("PORT", "3000"),
		MongoURI: getEnv("MONGO_URI", ""),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
		RedisTTL:  redisTTL,

		TimepadKey:   getEnv("TIMEPAD", ""),
		RouterKey:    getEnv("ROUTER_KEY", ""),
		RouterURL:    getEnv("ROUTER_URL", "https://api.openrouteservice.org"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),

	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
