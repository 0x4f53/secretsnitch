package main

import (
	"crypto/md5"
	"fmt"
	"strings"

	"golang.org/x/exp/rand"
)

func substringBeforeFirst(input string, delimiter string) string {
	index := strings.Index(input, delimiter)
	if index == -1 {
		return ""
	}
	return strings.TrimSpace(input[:index])
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func md5Hash(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func removeDuplicates(elements []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, element := range elements {
		if _, found := seen[element]; !found {
			seen[element] = struct{}{}
			result = append(result, element)
		}
	}

	return result
}
