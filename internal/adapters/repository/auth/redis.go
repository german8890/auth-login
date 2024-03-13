package auth

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // Select database 0
	})
}
