package store

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisURL string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
