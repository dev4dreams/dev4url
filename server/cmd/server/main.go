package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dev4dreams/dev4url/internal/config"
	"github.com/dev4dreams/dev4url/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type shortenUrl struct {
	ID          string `json:"id"`
	ShortenUrl  string `json:"shortenUrl"`
	OriginalUrl string `json:"originalUrl"`
	CreateTime  string `json:"createTime"`
}

var mockUrls = []shortenUrl{
	{
		ID:          "abc123",
		ShortenUrl:  "http://short.url/abc123",
		OriginalUrl: "https://www.example.com/very/long/original/url1",
		CreateTime:  "2025-01-13T14:30:00Z",
	},
	{
		ID:          "def456",
		ShortenUrl:  "http://short.url/def456",
		OriginalUrl: "https://www.example.com/another/long/original/url2",
		CreateTime:  "2025-01-13T14:35:00Z",
	},
	{
		ID:          "ghi789",
		ShortenUrl:  "http://short.url/ghi789",
		OriginalUrl: "https://www.example.com/yet/another/long/url3",
		CreateTime:  "2025-01-13T14:40:00Z",
	},
}

// GET request
func getUrlAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mockUrls)
}

// POST request
func createUrl(c *gin.Context) {
	var newUrl shortenUrl

	if err := c.BindJSON(&newUrl); err != nil {
		return
	}

	mockUrls = append(mockUrls, newUrl)
	c.IndentedJSON(http.StatusCreated, newUrl)
}

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.VerifyConnection(); err != nil {
		log.Fatalf("Connection verification failed: %v", err)
	}

	// customURL := "my-custom-url"
	// testURL := &db.URL{
	// 	ShortURL:    "abc123",
	// 	OriginalURL: "https://example.com",
	// 	CustomURL:   &customURL,
	// }

	// Insert the URL
	response, err := database.CreateURL(testURL)
	if err != nil {
		log.Fatalf("Failed to create URL: %v", err)
	}

	// Pretty print the response
	prettyResponse, _ := json.MarshalIndent(response, "", "  ")
	log.Printf("Successfully created URL:\n%s", string(prettyResponse))

	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	router.Use(cors.New(config))

	router.GET("/urlsAll", getUrlAll)
	router.POST("/createUrl", createUrl)

	router.Run("localhost:8080")
}
