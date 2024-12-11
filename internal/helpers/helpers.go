package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/weather-redis-api/internal/types"
)

func GetRedisValueByKey(key string) *types.WeatherResponse {

	var response types.WeatherResponse

	rdb := redis.NewClient(&redis.Options{
		Addr: "https://redis-production-0593.up.railway.app/:6379",
		Password: "hFDMGHeTNheyZWCZnzYVDbdNKytQsOJl",
		DB: 0,
	})

	ctx := context.Background()

	if val, err := rdb.Get(ctx, key).Result(); err == redis.Nil {
		fmt.Printf("Key %v does not exist or has expired.", key)
	} else if err != nil {
		fmt.Println("Error retrieving key:", err)
	} else {
		fmt.Printf("Retrieved key %s", key)
		json.Unmarshal([]byte(val), &response)
		return &response
	}

	return nil
}
