package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type News struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Url         string    `json:"url"`
		PublishedAt string    `json:"publishedAt"`
	} `json:"articles"`
}

var (
	NewsAPIKey string
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	NewsAPIKey = os.Getenv("NEWS_API_KEY")
	if NewsAPIKey == "" {
		log.Fatalf("NEWS_API_KEY must be set in the environment variables")
	}
}

func FetchNews() {
	fetchNewsWithKeyword("") // Fetch all top headlines
}

func FetchNewsByKeyword(keyword string) {
	fetchNewsWithKeyword(keyword) // Fetch news by specific keyword
}

func fetchNewsWithKeyword(keyword string) {
	baseUrl := "https://newsapi.org/v2/top-headlines?country=us"
	if keyword != "" {
		baseUrl += "&q=" + keyword
	}
	url := baseUrl + "&apiKey=" + NewsAPIKey

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatalf("Error fetching news: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err.Error())
	}

	var news News
	if err := json.Unmarshal(body, &news); err != nil {
		log.Fatalf("Error decoding news JSON: %s", err.Error())
	}

	for _, article := range news.Articles {
		log.Printf("News: %s - %s\n", article.Source.Name, article.Title)
	}
}

func main() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	go func() {
		for ; true; <-ticker.C {
			FetchNewsByKeyword("technology") // Example usage: Fetch news about technology
		}
	}()

	FetchNews()

	select {}
}