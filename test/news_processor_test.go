package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NewsData struct {
	Classification string
	Tags           []string
}

var (
	newsCache = make(map[string]*NewsData)
)

func ClassifyAndExtract(content string) (string, []string) {
	if data, found := newsCache[content]; found {
		return data.Classification, data.Tags
	}

	classification, tags := processContent(content)

	newsCache[content] = &NewsData{
		Classification: classification,
		Tags:           tags,
	}

	return classification, tags
}

func processContent(content string) (string, []string) {
	classification := "General"
	var tags []string

	if len(content) == 0 {
		classification = "Uncategorized"
	} else {
		if containsIgnoreCase(content, "technology") {
			classification = "Technology"
		} else if containsIgnoreCase(content, "finance") {
			classification = "Finance"
		}
		if containsIgnoreCase(content, "innovation") {
			tags = append(tags, "Innovation")
		}
		if containsIgnoreCase(content, "market") {
			tags = append(tags, "Market")
		}
	}

	return classification, tags
}

func containsIgnoreCase(content, keyword string) bool {
	return strings.Contains(strings.ToLower(content), strings.ToLower(keyword))
}

func TestNewsProcessing(t *testing.T) {
	classificationTestCases := []struct {
		content  string
		expected string
	}{
		{"The latest innovation in technology.", "Technology"},
		{"Market trends show a significant shift.", "General"},
		{"Something unrelated", "General"},
		{"", "Uncategorized"},
	}

	for _, tc := range classificationTestCases {
		classification, _ := ClassifyAndExtract(tc.content)
		assert.Equal(t, tc.expected, classification, "Expected and actual classification should match")
	}

	tagTestCases := []struct {
		content  string
		expected []string
	}{
		{"The latest innovation in technology.", []string{"Innovation"}},
		{"Market trends in the financial sector.", []string{"Market"}},
		{"A general news piece.", nil},
	}

	for _, tc := range tagTestCases {
		_, tags := ClassifyAndExtract(tc.content)
		assert.Equal(t, tc.expected, tags, "Expected and actual tags should match")
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