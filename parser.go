package main

import (
	"bufio"
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

func extractKeyValuePairs(text string) (map[string]string, error) {

	// Initialize a map to hold the key-value pairs
	keyValuePairs := make(map[string]string)

	// Regular expressions for matching various formats
	var reJSON = regexp.MustCompile(`"([^"]+)" *:\s*"([^"]+)"`)
	var reJS = regexp.MustCompile(`(\w+):"([^"]+)"`)
	var reEnv = regexp.MustCompile(`^\s*([A-Z_][A-Z0-9_]*?)\s*=\s*(.*)\s*$`)
	var reDict = regexp.MustCompile(`'([^']+)' *:\s*'([^']+)'`)
	var reGo = regexp.MustCompile(`\b(\w+)\s*:=\s*"([^"]+)"`)
	var reXML = regexp.MustCompile(`<([^\/>]+)[\/]*>.*</([^\/>]+)[\/]*>`)
	var reOtherLang = regexp.MustCompile(`\b(\w+)\s*=\s*"([^"]+)"`)

	// Scan the file line by line
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()

		// Match JSON style key-value pairs
		if matches := reJSON.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}

		// Match JS style key-value pairs
		if matches := reJS.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}

		// Match .env file style key-value pairs
		if matches := reEnv.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}

		// Match Dict style key-value pairs
		if matches := reDict.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}

		// Match Go language key-value pairs
		if matches := reGo.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}

		// Match XML key-value pairs
		if matches := reXML.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				value := strings.Replace(match[0], match[1], "", -1)
				value = strings.Replace(value, "<>", "", -1)
				value = strings.Replace(value, "</>", "", -1)
				keyValuePairs[match[1]] = value
			}
		}

		// Match other programming languages key-value pairs
		if matches := reOtherLang.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				keyValuePairs[match[1]] = match[2]
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keyValuePairs, nil
}
