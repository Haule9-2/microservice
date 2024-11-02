package main

import (
    "context"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Change to "redis:6379" if using Docker
    })
    defer client.Close() // Close the client when done

    // Test Redis connection
    err := client.Set(ctx, "test_key", "test_value", 0).Err()
    if err != nil {
        log.Fatalf("Failed to set key: %v", err)
    }

    val, err := client.Get(ctx, "test_key").Result()
    if err != nil {
        log.Fatalf("Failed to get key: %v", err)
    }

    log.Printf("Value for 'test_key': %s", val)
}
