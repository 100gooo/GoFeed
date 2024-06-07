package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type UserPreferences struct {
	FavoriteCategories []string `json:"favoriteCategories"`
	Notification bool `json:"notification"`
	FeedPresentation string `json:"feedPresentation"`
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
	_, exists := p.preferences[userID]
	if !exists {
		return errors.New("user not found")
	}

	p.preferences[userID] = preferences
	return nil
}

func (p *PreferencesService) GetPreferences(userID string) (UserPreferences, error) {
	preferences, exists := p.preferences[userID]
	if !exists {
		return UserPreferences{}, errors.New("user not found")
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
		Notification: true,
		FeedPresentation: "Compact",
	}

	err := service.UpdatePreferences(userID, pref)
	if err != nil {
		fmt.Println("Error updating preferences:", err)
		return
	}

	updatedPref, err := service.GetPreferences(userID)
	if err != nil {
		fmt.Println("Error fetching preferences:", err)
		return
	}

	prefJson, _ := json.Marshal(updatedPref)
	fmt.Println("Updated Preferences:", string(prefJson))
}