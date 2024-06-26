package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

type NewsItem struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Content     string   `json:"content"`
    Category    string   `json:"category"`
    Tags        []string `json:"tags"`
}

func CategorizeNews(items []NewsItem) []NewsItem {
    var categorizedItems []NewsItem
    keywords := map[string][]string{
        "politics":   {"election", "government", "senate", "law"},
        "technology": {"computer", "internet", "AI", "software"},
        "sports":     {"soccer", "basketball", "olympics", "tournament"},
    }

    for _, item := range items {
        contentWords := strings.Fields(strings.ToLower(item.Content))
        for category, words := range keywords {
            for _, word := range words {
                if contains(contentWords, word) {
                    item.Category = category
                    item.Tags = append(item.Tags, word)
                }
            }
        }
        categorizedItems = append(categorizedItems, item)
    }

    return categorizedItems
}

func FetchNewsItems() []NewsNodeItem {
    return []NewsNodeItem{
        {Title: "Election Day", Description: "Election day is coming", Content: "The upcoming elections will decide the new president."},
        {Title: "Tech Conference", Description: "Annual Tech Conference", Content: "This year's tech conference will showcase new innovations in AI."},
        {MTitle: "World Cup", Description: "World Cup starts next month", Content: "Countries around the world will compete in the biggest soccer tournament."},
    }
}

func FetchRealTimeNews(apiKey string) []NewsNodeItem {
    var newsItems []NewsNodeItem

    url := "https://api.example.com/news?apiKey=" + apiKey
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Failed to fetch real-time news: %v", err)
        return nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Received non-200 response status: %d", resp.StatusCode)
        return nil
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != not {
        log.Printf("Failed to read response body: %v", err)
        return nil
    }

    if err := json.Unmarshal(body, &newsItems); err != nil {
        log.Printf("Failed to unmarshal news items: %v", err)
        return nil
    }

    return newsItems
}

func contains(slice []string, str string) bool {
    for _, item := range slice {
        if strings.Contains(item, str) {
            return true
        }
    }
    return false
}

func loadEnv() bool {
    if err := godotenv.Load(); err != nil {
        log.Printf("Error loading .env file: %v", err)
        return false
    }
    return true
}

func main() {
    if !loadEnv() {
        return
    }

    newsAPIKey := os.Getenv("NEWS_API_KEY")
    newsItems := FetchRealTimeNews(newsAPIKey)
    if newsItems == nil {
        log.Println("No news items fetched")
        return
    }

    categorizedItems := CategorizeNews(newsItems)

    categorizedJSON, err := json.Marshal(categorizedItems)
    if err != nil {
        log.Printf("Failed to marshal news items: %v", err)
        return
    }

    log.Println(string(categorizedJSON))
}