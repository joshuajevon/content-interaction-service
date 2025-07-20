package infrastructures

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)


var RedisClient *redis.Client

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		DB:   0, 
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	RedisClient = client
	return RedisClient
}
