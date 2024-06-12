package news_fetcher_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"yourapp/news_fetcher"
)

func TestFetchNewsSuccessfully(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","articles":[{"title":"Fake News Title","description":"Fake News Description"}]}`))
	}))
	defer mockAccServer.Close()

	os.Setenv("NEWS_API_URL", mockServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	newsItems, err := news_fetcher.FetchNews()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(newsItems) != 1 {
		t.Fatalf("Expected 1 news item, got %d", len(newsItems))
	}

	if newsItems[0].Title != "Fake News Title" {
		t.Errorf("Expected title to be 'Fake News Title', got '%s'", newsItems[0].Title)
	}
}

func TestFetchNewsWithError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	os.Setenv("NEWS_API_URL", mockServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	_, err := news_fetcher.FetchNews()
	if err == nil {
		t.Fatal("Expected an error, got none")
	}

	var apiErr *news_fetcher.APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected an APIError, got %T", err)
	}
}

func TestFetchNewsWithInvalidJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer mockServer.Close()

	os.Setenv("NEWS_API_URL", mockServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	_, err := news_fetcher.FetchNews()
	if err == nil {
		t.Fatal("Expected an error, got none")
	}

	if err != nil && !errors.As(err, &json.SyntaxError{}) {
		t.Errorf("Expected a JSON syntax error, got %v", err)
	}
}