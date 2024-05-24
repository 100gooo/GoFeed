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

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Source      string `json:"source"`
}

var articles []Article

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	articles = []Article{
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

func fetchArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func fetchArticlesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category := strings.ToLower(r.URL.Query().Get("category"))

	var filteredArticles []Article
	for _, article := range articles {
		if strings.ToLower(article.Category) == category {
			filteredArticles = append(filteredArticles, article)
		}
	}

	json.NewEncoder(w).Encode(filteredArticles)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/articles", fetchArticles).Methods("GET")
	router.HandleFunc("/articles/search", fetchArticlesByCategory).Methods("GET")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error occurred while running the server: %s", err)
	}
}