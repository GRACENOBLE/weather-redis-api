package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/weather-redis-api/internal/helpers"
	"github.com/weather-redis-api/internal/types"
)

func WeatherServer() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults.")
	}

	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	r.GET("/weather", func(c *gin.Context) {
		lat := c.Query("lat") //"0.3976677749854413"
		lon := c.Query("lon") //"32.6378629998115"
		key := lat + lon

		dbData := helpers.GetRedisValueByKey(key)

		if dbData == nil {

			if lat == "" || lon == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
				return
			}

			apiKey := os.Getenv("OPENWEATHER_API_KEY")
			if apiKey == "" {
				log.Fatal("OPENWEATHER_API_KEY is not set in .env")
			}

			url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

			fmt.Printf("\nurl: %v\n", url)

			res, err := http.Get(url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather"})
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				c.JSON(res.StatusCode, gin.H{"error": "Invalid request"})
				return
			}

			var weather types.WeatherResponse

			if bodyBytes, err := io.ReadAll(res.Body); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			} else {

				if err := rdb.Set(ctx, key, bodyBytes, 30*time.Minute).Err(); err != nil {
					log.Fatalf("Error setting key:%v", err)
				}

				if err := json.Unmarshal(bodyBytes, &weather); err != nil {
					log.Printf("%v",err.Error())
				}

				fmt.Printf("Key set successfully!")

				c.JSON(200, weather)

				return
			}
		}

		c.JSON(200, dbData)

	})

	r.GET("/database", func(c *gin.Context) {
		dbData := helpers.GetRedisValueByKey("weatherToday")
		fmt.Printf("%v", dbData)
		c.JSON(200, dbData)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	r.Run(":" + port) // listen and serve on port 0.0.0.0:8081 or 0.0.0.0:8080 by default
}
