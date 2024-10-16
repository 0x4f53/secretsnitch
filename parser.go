package main

import (
	"regexp"
	"strings"
)

func containsBlacklisted(text string) bool {

	var blacklist = []string{
		"data:image/png;",
	}

	for _, item := range blacklist {
		if strings.Contains(text, item) {
			return true
		}
	}

	return false

}

func grabDeclaredStringValues(text string) []string {

	var capturedStrings []string

	stringPattern := `("([^"\\]*(\\.[^"\\]*)*)")|('([^'\\]*(\\.[^'\\]*)*)')|(` + "`([^`\\\\]*(\\\\.[^`\\\\]*)*)`" + `)`
	re := regexp.MustCompile(stringPattern)

	matches := re.FindAllString(text, -1)
	for _, match := range matches {
		if !containsBlacklisted(match) {
			capturedStrings = append(capturedStrings, match)
		}
	}

	return capturedStrings
}
