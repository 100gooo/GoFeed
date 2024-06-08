package main

import (
	"encoding/json"
	"fmt"
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

func (p *PreferencesService) GetPreferences(userID string) (UserPreferences, error) {
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

func main() {
	LoadConfiguration()

	service := NewPreferencesService()

	userID := "123"
	pref := UserPreferences{
		FavoriteCategories: []string{"Tech", "Sports"},
		Notification:       true,
		FeedPresentation:   "Compact",
	}
	service.preferences[userID] = pref

	err := service.UpdatePreferences(userID, pref)
	if err != nil {
		if _, ok := err.(*UserNotFoundError); ok {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Unexpected error updating preferences:", err)
		}
		return
	}

	updatedPref, err := service.GetPreferences(userID)
	if err != nil {
		fmt.Println("Error fetching preferences:", err)
		return
	}
	prefJson, err := json.Marshal(updatedPref)
	if err != nil {
		fmt.Println("Error marshalling preferences to JSON:", err)
		return
	}
	fmt.Println("Updated Preferences:", string(prefJson))
}