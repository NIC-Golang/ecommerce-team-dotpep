package redis_test

import (
	"testing"

	"context"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisSetGet(t *testing.T) {
	mockRedis := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mockRedis.Addr()})
	ctx := context.Background()

	client.Set(ctx, "key", "value", 0)
	val, err := client.Get(ctx, "key").Result()

	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}
