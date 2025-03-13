package redis

import (
	"cart-service/golang/internal/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bytedance/sonic"

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

func getCartKey(id string) string {
	return "cart:" + id
}

func GetCartFromRedis(id string) (*models.Cart, error) {
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
	if err := sonic.Unmarshal(res, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func SaveToCart(id string, cart *models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	jsonCart, err := sonic.Marshal(cart)
	if err != nil {
		return err
	}

	err = GetRedisClient().Set(ctx, getCartKey(id), jsonCart, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteCartFromRedis(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleted, err := client.Del(ctx, getCartKey(id)).Result()
	if err != nil {
		return fmt.Errorf("error deleting cart: %w", err)
	}

	if deleted == 0 {
		return nil
	}

	return nil
}

func FindCartItem(productId, id string) (*models.CartItem, error) {
	cart, err := GetCartFromRedis(id)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}
	for _, item := range cart.Items {
		if item.ProductID == productId {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item not found in cart")
}
