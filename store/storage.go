package store

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedisStorage() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Redis: %v", err)
	}

	return client, nil
}
