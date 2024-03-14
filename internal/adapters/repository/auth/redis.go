package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1, // Select database 1
	})

	// Ping the Redis server to check if the connection is successful
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis server: %v", err)
	}

	fmt.Println("Redis server responded:", pong)

	return client, nil
}
