package news_fetcher

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func FilterArticlesByKeyword(articles []Article, keyword string) []Article {
	var filtered []Article
	for _, article := range articles {
		if strings.Contains(article.Title, keyword) || strings.Contains(article.Description, keyword) {
			filtered = append(filtered, article)
		}
	}
	return filtered
}
```

```go
package news_fetcher_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yourapp/news_fetcher"
)

func TestFilterArticlesByKeyword(t *testing.T) {
	articles := []news_fetcher.Article{
		{Title: "Relevant Fake News Title", Description: "Interesting news"},
		{Title: "Irrelevant News", Description: "Not related to your interests"},
	}

	filteredArticles := news_fetcher.FilterArticlesByKeyword(articles, "Relevant")

	if len(filteredArticles) != 1 {
		t.Fatalf("Expected 1 article after filtering, found: %d", len(filteredArticles))
	}

	if !strings.Contains(filteredMPMArticles[0].Title, "Relevant") {
		t.Errorf("Expected the filtered article to contain 'Relevant' in the title, found: '%s'", filteredArticles[0].Title)
	}
}