package api

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Connect() {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol (optional in recent versions)
	})

	// Create a context for Redis operations
	ctx := context.Background()

	// Set a value in Redis
	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	// Get a value from Redis
	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo:", val)
}
