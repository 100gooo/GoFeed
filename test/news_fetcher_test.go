package news_fetcher_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"yourapp/news_fetcher"
)

func TestFetchNewsSucceeds(t *testing.T) {
	mockNewsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","articles":[{"title":"Fake News Title","description":"Fake News Description"}]}`))
	}))
	defer mockNewsServer.Close()

	os.Setenv("NEWS_API_URL", mockNewsServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	fetchedArticles, err := news_fetcher.FetchNews()
	if err != nil {
		t.Fatalf("Did not expect an error, received: %v", err)
	}

	if len(fetchedArticles) != 1 {
		t.Fatalf("Expected 1 article, found: %d", len(fetchedArticles))
	}

	if fetchedArticles[0].Title != "Fake News Title" {
		t.Errorf("Expected article title to be 'Fake News Title', found: '%s'", fetchedArticles[0].Title)
	}
}

func TestFetchNewsReturnsErrorOnServerFailure(t *testing.T) {
	mockErrorResponseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockErrorResponseServer.Close()

	os.Setenv("NEWS_API_URL", mockErrorResponseServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	_, err := news_fetcher.FetchNews()
	if err == nil {
		t.Fatal("Expected an error, received none")
	}

	var fetchError *news_fetcher.APIError
	if !errors.As(err, &fetchError) {
		t.Errorf("Expected an APIError, received: %T", err)
	}
}

func TestFetchNewsHandlesInvalidJSONResponse(t *testing.T) {
	mockInvalidJsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockInvalidJsonServer.Close()

	os.Setenv("NEWS_API_URL", mockInvalidJsonServer.URL)
	defer os.Unsetenv("NEWS_API_URL")

	_, err := news_fetcher.FetchNews()
	if err == nil {
		t.Fatal("Expected an error, received none")
	}

	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("Expected a JSON syntax error, received: %v", err)
	}
}