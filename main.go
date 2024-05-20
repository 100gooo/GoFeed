package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type News struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Source      string   `json:"source"`
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
	}
}

func getNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsItems)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/getNews", getNews).Methods("GET")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error occurred while running the server: %s", err)
	}
}