package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RedisConnection(address string) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %s", err)
	}

	fmt.Println("Connected to Redis")
	return client, nil
}
