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

func addArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newArticle Article
	if err := json.NewDecoder(r.Body).Decode(&newFromRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid article format")
		return
	}

	// Basic Validation
	if newArticle.Title == "" || newArticle.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing title or description")
		return
	}

	articles = append(articles, newArticle)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newArticle)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/articles", fetchArticles).Methods("GET")
	router.HandleFunc("/articles/search", fetchArticlesByCategory).Methods("GET")
	router.HandleFunc("/articles", addArticle).Methods("POST")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error occurred while running the server: %s", err)
	}
}