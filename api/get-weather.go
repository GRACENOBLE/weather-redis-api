package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weather-redis-api/internal/types"
)

func Gin() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults.")
	}

	r := gin.Default()
	r.GET("/weather", func(c *gin.Context) {
		lat := c.Query("lat") //"0.3976677749854413"
		lon := c.Query("lon") //"32.6378629998115"

		if lat == "" || lon == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
			return
		}
		apiKey := os.Getenv("OPENWEATHER_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENWEATHER_API_KEY is not set in .env")
		}
		fmt.Println(apiKey)

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
			err := json.Unmarshal(bodyBytes, &weather)
			if err != nil {
				log.Printf(err.Error())
			}

			c.JSON(200, weather)
		}

	})
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	r.Run(":" + port) // listen and serve on port 0.0.0.0:8081 or 0.0.0.0:8080 by default
}
