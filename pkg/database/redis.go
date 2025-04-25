package database

import (
	"context"
	"fmt"
	"product_recommendation/pkg/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var once sync.Once

type RedisClientInterface interface {
	GetRedisClient() *redis.Client
	RedisPing(ctx context.Context) error
}

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient() RedisClientInterface {
	return &RedisClient{}
}

func (r *RedisClient) GetRedisClient() *redis.Client {
	once.Do(func() {
		config := config.GetConfig()
		r.redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})
	})
	return redisClient
}

func (r *RedisClient) RedisPing(ctx context.Context) error {
	return r.redisClient.Ping(ctx).Err()
}
