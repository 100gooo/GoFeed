package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	newsClassificationCache = make(map[string]string)
	newsTagsCache           = make(map[string][]string)
)

func ClassifyNewsContent(content string) string {
	if classification, found := newsClassificationCache[content]; found {
		return classification
	}

	classification := "General"
	if len(content) == 0 {
		classification = "Uncategorized"
	} else if containsIgnoreCase(content, "technology") {
		classification = "Technology"
	} else if containsIgnoreCase(content, "finance") {
		classification = "Finance"
	}

	newsClassificationCache[content] = classification
	return classification
}

func ExtractTagsFromContent(content string) []string {
	if tags, found := newsTagsCache[content]; found {
		return tags
	}

	var tags []string
	if containsIgnoreCase(content, "innovation") {
		tags = append(tags, "Innovation")
	}
	if containsIgnoreCase(content, "market") {
		tags = append(tags, "Market")
	}

	newsTagsCache[content] = tags
	return tags
}

func containsIgnoreCase(content, keyword string) bool {
	return strings.Contains(strings.ToLower(content), strings.ToLower(keyword))
}

func TestNewsClassification(t *testing.T) {
	testCases := []struct {
		content  string
		expected string
	}{
		{"The latest innovation in technology.", "Technology"},
		{"Market trends show a significant shift.", "General"},
		{"Something unrelated", "General"},
		{"", "Uncategorized"},
	}

	for _, tc := range testCases {
		actualClassification := ClassifyNewsContent(tc.content)
		assert.Equal(t, tc.expected, actualClassification, "Expected and actual classification should match")
	}
}

func TestNewsTagExtraction(t *testing.T) {
	testCases := []struct {
		content  string
		expected []string
	}{
		{"The latest innovation in technology.", []string{"Innovation"}},
		{"Market trends in the financial sector.", []string{"Market"}},
		{"A general news piece.", nil},
	}

	for _, tc := range testCases {
		actualTags := ExtractTagsFromContent(tc.content)
		assert.Equal(t, tc.expected, actualTags, "Expected and actual tags should match")
	}
}

func SetupEnvironmentVariables() {
	os.Setenv("API_KEY", "YourAPIKeyHere")
}

func TestMain(m *testing.M) {
	SetupEnvironmentVariables()
	code := m.Run()
	os.Exit(code)
}