package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Source      string `json:"source"`
}

var newsItems []News

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	newsItems = []News{
		{
			Title:       "Title 1",
			Description: "Description 1",
			Category:    "Technology",
			Source:      "Source 1",
		},
		{
			Title:       "Title 2",
			Description: "Description 2",
			Category:    "Science",
			Source:      "Source 2",
		},
	}
}

func getNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsItems)
}

func getNewsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category := r.URL.Query().Get("category")

	var filteredNews []News
	for _, item := range newsItems {
		if strings.EqualFold(item.Category, category) {
			filteredNews = append(filteredNews, item)
		}
	}

	json.NewEncoder(w).Encode(filteredNews)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/getNews", getNews).Methods("GET")
	r.HandleFunc("/searchNews", getNewsByCategory).Methods("GET")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error occurred while running the server: %s", err)
	}
}