package cacheclient

import (
    "context"
    "github.com/go-redis/redis/v8"
    "log"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })
    _, err := client.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
    return client
}
