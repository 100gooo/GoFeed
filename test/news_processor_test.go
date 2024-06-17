package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockClassifyNews(content string) string {
	if len(content) == 0 {
		return "Uncategorized"
	}
	if contains(content, "technology") {
		return "Technology"
	}
	if contains(content, "finance") {
		return "Finance"
	}
	return "General"
}

func mockTagMetadata(content string) []string {
	tags := []string{}
	if contains(content, "innovation") {
		tags = append(tags, "Innovation")
	}
	if contains(content, "market") {
		tags = append(tags, "Market")
	}
	return tags
}

func contains(content, keyword string) bool {
	return strings.Contains(strings.ToLower(content), strings.ToLower(keyword))
}

func TestClassifyNews(t *testing.T) {
	testCases := []struct {
		content  string
		expected string
	}{
		{"The latest innovation in technology.", "Technology"},
		{"Market trends show a significant shift.", "Finance"},
		{"Something unrelated", "General"},
		{"", "Uncategorized"},
	}

	for _, tc := range testCases {
		result := mockClassifyNews(tc.content)
		assert.Equal(t, tc.expected, result, "They should be equal")
	}
}

func TestTagMetadata(t *testing.T) {
	testCases := []struct {
		content  string
		expected []string
	}{
		{"The latest innovation in technology.", []string{"Innovation"}},
		{"Market trends in the financial sector.", []string{"Market"}},
		{"A general news piece.", nil},
	}

	for _, tc := range testCases {
		result := mockTagMetadata(tc.content)
		assert.Equal(t, tc.expected, result, "The tags should match")
	}
}

func setUpEnv() {
	os.Setenv("API_KEY", "YourAPIKeyHere")
}

func TestMain(m *testing.M) {
	setUpEnv()
	code := m.Run()
	os.Exit(code)
}