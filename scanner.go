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
	"mvdan.cc/xurls/v2"
)

type Secret struct {
	Tool            string
	ScanTimestamp   string
	Provider        string
	ServiceName     string
	Matches         []string
	Entropy         float64
	CacheFile       string
	URL             string
	CapturedDomains []string
	CapturedURLs    []string
	Tags            []string
}

func getMatchingLines(input string, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var matches []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		matches = append(matches, re.FindStringSubmatch(line)...)
	}

	return matches, nil

}

func grabURLs(text string) []string {

	var captured []string
	location := substringBeforeFirst(text, "---")

	scanner := bufio.NewScanner(strings.NewReader(text))

	rx := xurls.Relaxed()

	for scanner.Scan() {
		line := scanner.Text()
		urls := rx.FindAllString(line, -1)
		for _, url := range urls {
			if strings.Contains(url, "://") && url != location {
				captured = append(captured, url)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading string: %s\n", err)
	}

	return captured

}

func FindSecrets(text string) Secret {

	var secret Secret

	var tags []string

	domains, _ := textsubs.DomainsOnly(text, false)
	domains = textsubs.Resolve(domains)

	for _, service := range signatures {
		for keyName, regex := range service.Keys {
			matches, err := getMatchingLines(text, regex)
			if err != nil {
				//log.Printf("Error reading data: %v\n", err)
				return secret
			}

			if len(matches) > 0 {

				tags = append(tags, "regexMatch")
				url := substringBeforeFirst(text, "---")
				capturedUrls := grabURLs(text)

				entropy := EntropyPercentage(text)
				if entropy > 60 {
					tags = append(tags, "entropic")
				}

				secret = Secret{
					Tool:            "secretsnitch",
					ScanTimestamp:   time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
					Provider:        service.Name,
					ServiceName:     keyName,
					Entropy:         entropy,
					URL:             url,
					Matches:         matches,
					CapturedDomains: domains,
					CapturedURLs:    removeDuplicates(capturedUrls),
					Tags:            removeDuplicates(tags),
				}
			}
		}
	}

	return secret

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

	if len(secrets.Matches) > 0 {
		unindented, _ := json.Marshal(secrets)
		appendToFile(*outputFile, string(unindented))
		indented, _ := json.MarshalIndent(secrets, "", "	")
		fmt.Println(string(indented))
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
