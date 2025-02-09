package redis

import (
	"cart-service/golang/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
)

func InitRedis(address string) error {
	var err error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
			Protocol: 2,
		})
		_, err = client.Ping(context.Background()).Result()
	})

	if err != nil {
		log.Printf("Redis connection error: %v", err)
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func GetRedisClient() *redis.Client {
	if client == nil {
		panic("Redis client is not initialized. Call InitRedis() first.")
	}
	return client
}

func getCartKey(id int) string {
	return "cart:" + strconv.Itoa(id)
}

func GetCartFromRedis(id int) (*models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := GetRedisClient()

	res, err := client.Get(ctx, getCartKey(id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal(res, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func SaveToCart(id int, cart *models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	jsonCart, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	err = client.Set(ctx, getCartKey(id), jsonCart, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
