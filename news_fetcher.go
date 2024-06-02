package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

type NewsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

var (
	newsAPIKey string
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	newsAPIKey = os.Getenv("NEWS_API_KEY")
	if newsAPIKey == "" {
		log.Fatalf("NEWS_API_KEY must be set in the environment variables")
	}
}

func init() {
	loadConfig()
}

func fetchLatestHeadlines() {
	fetchNews("") // Fetch all top headlines without a keyword
}

func fetchNewsByKeyword(keyword string) {
	fetchNews(keyword) // Fetch headlines by a specific keyword
}

func fetchNews(keyword string) {
	newsAPIBaseURL := "https://newsapi.org/v2/top-headlines?country=us"
	if keyword != "" {
		newsAPIBaseURL += "&q=" + keyword
	}
	requestURL := newsAPIBasefluxURL + "&apiKey=" + newsAPIKey

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		log.Printf("Error fetching news: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API request failed with status code %d: %s", resp.StatusCode, resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	var news NewsResponse
	if err := json.Unmarshal(body, &news); err != nil {
		log.Printf("Error decoding news JSON: %v", err)
		return
	}

	if news.Status != "ok" {
		log.Printf("Error with the news API response: %s", news.Status)
		return
	}

	for _, article := range news.Articles {
		log.Printf("Article: %s - %s\n", article.Source.Name, article.Title)
	}
}

func main() {
	newsUpdateTicker := time.NewTicker(30 * time.Minute)
	defer newsUpdateTicker.Stop()

	go func() {
		for range newsUpdateTicker.C {
			fetchNewsByKeyword("technology") // Example usage: Fetch news about technology
		}
	}()

	fetchLatestHeadlines()
	select {}
}