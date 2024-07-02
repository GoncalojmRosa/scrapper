package store

import (
	"github.com/go-redis/redis"
)

func NewRedisStorage(conString string) (*redis.Client, error) {
	opt, err := redis.ParseURL(conString)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(opt)

	return redisClient, nil
}
