package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/0x4f53/textsubs"
	"github.com/dlclark/regexp2"
	"mvdan.cc/xurls/v2"
)

type Secret struct {
	Provider    string
	ServiceName string
	Variable    string
	Secret      string
	Entropy     float64
	Tags        []string
}

type ToolData struct {
	Tool            string
	ScanTimestamp   string
	Secrets         []Secret
	CacheFile       string
	SourceUrl       string
	CapturedDomains []string
	CapturedURLs    []string
}

func getMatchingLines(input string, pattern string) (map[string]string, error) {

	// This is not the regular golang library. It supports lookaheads and stuff.
	re := regexp2.MustCompile(pattern, 0)

	matches := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		values, _ := extractKeyValuePairs(line)
		for key, value := range values {
			match, _ := re.MatchString(value)
			if match && !containsBlacklisted(value) {
				matches[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil && len(matches) == 0 {
		values, _ := extractKeyValuePairs(input)
		for key, value := range values {
			match, _ := re.MatchString(value)
			if match {
				matches[key] = value
			}
		}
	}

	return matches, nil

}

func grabURLs(text string) []string {

	var captured []string
	location := substringBeforeFirst(text, "---")

	scanner := bufio.NewScanner(strings.NewReader(text))

	rx := xurls.Relaxed()
	rxUrls := rx.FindAllString(text, -1)
	captured = append(captured, rxUrls...)

	splitText := strings.Split(text, "{")

	protocol := substringBeforeFirst(location, "://")

	for _, line := range splitText {

		re := regexp.MustCompile(`(?:href|src|action|cite|data|formaction|poster)\s*=\s*["']([^"']+)["']`)
		matches := re.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			fixedUrl := match[1]
			if strings.HasPrefix(fixedUrl, "//") {
				fixedUrl = protocol + ":" + fixedUrl
			}
			validDomains, _ := textsubs.DomainsOnly(fixedUrl, false)
			if len(validDomains) > 0 && !strings.Contains(fixedUrl, "://") {
				fixedUrl = protocol + ":" + fixedUrl
			}
			captured = append(captured, fixedUrl)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading string: %s\n", err)
	}

	var urls []string
	for _, url := range captured {
		if strings.Contains(url, "://") {
			urls = append(urls, url)
		}
	}

	return removeDuplicates(urls)

}

func FindSecrets(text string) ToolData {

	var output ToolData
	var secrets []Secret

	var tags []string

	domains, _ := textsubs.DomainsOnly(text, false)
	domains = textsubs.Resolve(domains)

	splitText := strings.Split(text, "{")

	var mu sync.Mutex
	var wg sync.WaitGroup

	workerCount := 10000
	lineChan := make(chan string, workerCount)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				data, _ := extractKeyValuePairs(line)

				for key, value := range data {
					for _, provider := range signatures {
						for service, regex := range provider.Keys {
							re := regexp2.MustCompile(regex, 0)
							match, _ := re.MatchString(value)

							if match {
								mu.Lock()
								tags = append(tags, "regexMatched")
								mu.Unlock()

								entropy := EntropyPercentage(value)
								if entropy > 66.6 {
									mu.Lock()
									tags = append(tags, "highEntropy")
									mu.Unlock()
								}

								providerString := strings.ToLower(strings.Split(provider.Name, ".")[0])
								if strings.Contains(strings.ToLower(text), providerString) {
									mu.Lock()
									tags = append(tags, "providerDetected")
									mu.Unlock()
								}

								secret := Secret{
									Provider:    provider.Name,
									ServiceName: service,
									Variable:    key,
									Secret:      value,
									Entropy:     entropy,
									Tags:        removeDuplicates(tags),
								}

								mu.Lock()
								secrets = append(secrets, secret)
								mu.Unlock()
							}
						}
					}
				}
			}
		}()
	}

	for _, line := range splitText {
		lineChan <- line
	}

	close(lineChan)

	wg.Wait()

	sourceUrl := substringBeforeFirst(text, "---")
	capturedUrls := grabURLs(text)

	output = ToolData{
		Tool:            "secretsnitch",
		ScanTimestamp:   time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		SourceUrl:       sourceUrl,
		Secrets:         secrets,
		CapturedDomains: domains,
		CapturedURLs:    removeDuplicates(capturedUrls),
	}

	return output
}

func scanFile(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	secrets := FindSecrets(string(data))
	secrets.CacheFile = filePath

	if len(secrets.Secrets) > 0 {
		unindented, _ := json.Marshal(secrets)
		appendToFile(*outputFile, string(unindented))
		indented, _ := json.MarshalIndent(secrets, "", "	")
		fmt.Println(string(indented))
	}

	if *recurse {
		urls := grabURLs(string(data))
		fetchFromUrlList(urls)
		files, _ := listCachedFiles()
		ScanFiles(files)
	}

}

func ScanFiles(files []string) {

	var wg sync.WaitGroup
	fileChan := make(chan string)

	for i := 0; i < maxWorkers; i++ {
		go scanWorker(fileChan, &wg)
	}

	for _, file := range files {
		wg.Add(1)
		fileChan <- file
	}

	close(fileChan)
	wg.Wait()
}

func scanWorker(files <-chan string, wg *sync.WaitGroup) {
	for file := range files {
		scanFile(file, wg)
	}
}
