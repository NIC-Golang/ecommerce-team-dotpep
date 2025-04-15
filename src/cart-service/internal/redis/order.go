package redis

import (
	"cart-service/golang/internal/models"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

func CreateOrder(id string, cart *models.Cart, createdAT time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeDuration)
	defer cancel()
	total := 0.0
	for _, val := range cart.Items {
		total += val.Price
	}
	order := &models.Order{
		OrderNumber: generateOrderId(),
		UserID:      cart.UserID,
		Status:      "confirmed",
		Items:       cart.Items,
		TotalPrice:  total,
		CreatedAt:   createdAT,
	}
	orderJSON, err := sonic.Marshal(order)
	if err != nil {
		return err
	}
	err = getRedisClient().Set(ctx, getOrderKey(id), orderJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func generateOrderId() string {
	return fmt.Sprintf("ORDER-%06d", rand.Intn(1000000))
}

func GetOrderFromRedis(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeDuration)
	defer cancel()

	result, err := getRedisClient().Get(ctx, getOrderKey(id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var order *models.Order
	err = sonic.Unmarshal(result, &order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func SaveToOrder(id string, order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeDuration)
	defer cancel()
	orderJSON, err := sonic.Marshal(order)
	if err != nil {
		return err
	}
	err = getRedisClient().Set(ctx, getOrderKey(id), orderJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
