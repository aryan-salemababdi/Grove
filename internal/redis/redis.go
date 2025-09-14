package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

func InitRedis(addr, password string, db int, container *dig.Container) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Redis connected successfully")

	if err := container.Provide(func() *redis.Client { return rdb }); err != nil {
		log.Fatalf("❌ Failed to provide Redis client to container: %v", err)
	}

	return rdb
}
