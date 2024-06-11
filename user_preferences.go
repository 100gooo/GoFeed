package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserNotFoundError struct {
	UserID string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with ID %s not found", e.UserID)
}

type UserPreferences struct {
	FavoriteCategories []string `json:"favoriteCategories"`
	Notification       bool     `json:"notification"`
	FeedPresentation   string   `json:"feedPresentation"`
}

type PreferencesService struct {
	preferences map[string]UserPreferences
}

func NewPreferencesService() *PreferencesService {
	return &PreferencesService{
		preferences: make(map[string]UserPreferences),
	}
}

func (p *PreferencesService) UpdatePreferences(userID string, preferences UserPreferences) error {
	if _, exists := p.preferences[userID]; !exists {
		return &UserNotFoundError{UserID: userID}
	}
	p.preferences[userID] = preferences
	return nil
}

func (p *PreferenceseService) GetPreferences(userID string) (UserPreferences, error) {
	preferences, exists := p.preferences[userID]
	if !exists {
		return UserPreferences{}, &UserNotFoundError{UserID: userID}
	}
	return preferences, nil
}

func LoadConfiguration() error {
	os.Setenv("NOTIFICATION_DEFAULT", "true")
	return nil
}

// HTTP Handlers
func (p *PreferencesService) updatePreferencesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	var userID = r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	var preferences UserPreferences
	err := json.NewDecoder(r.Body).Decode(&preferences)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = p.UpdatePreferences(userID, preferences)
	if err != nil {
		if _, ok := err.(*UserNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Unexpected error updating preferences", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *PreferencesService) getPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	var userID = r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	preferences, err := p.GetPreferences(userID)
	if err != nil {
		if _, ok := err.(*UserNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Unexpected error fetching preferences", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(preferences)
}

func (p *PreferencesService) registerHandlers() {
	http.HandleFunc("/preferences", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			p.getPreferencesHandler(w, r)
		case "POST":
			p.updatePreferencesHandler(w, r)
		default:
			http.Error(w, "Method not supported", http.StatusNotFound)
		}
	})
}

func main() {
	LoadConfiguration()

	service := NewPreferencesService()

	// Start HTTP server
	service.registerHandlers()
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}