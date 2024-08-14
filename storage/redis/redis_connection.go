package redis

import (
	"budgeting-service/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis_HOST, cfg.Redis_PORT),
		Password: fmt.Sprint(cfg.Redis_PASSWORD),
		DB:       cfg.Redis_DB,
	})

	return client
}
