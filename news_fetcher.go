package main

import (
	"encoding/json"
	"io/ioutil"
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

var NewsAPIKey string

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
	url := "https://newsapi.org/v2/top-headlines?country=us&apiKey=" + NewsAPIKey

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching news: %s", err.Error())
	}
	defer resp.Body.Close()

	var news News
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
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
			FetchNews()
		}
	}()

	FetchNews()

	select {}
}