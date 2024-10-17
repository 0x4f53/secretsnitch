package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func getMatchingLines(input string, pattern string) ([]string, error) {

	// This is not the regular golang library. It supports lookaheads and stuff.
	re := regexp2.MustCompile(pattern, 0)

	var matches []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		values := grabDeclaredStringValues(line)
		for _, value := range values {
			match, _ := re.MatchString(value)
			if match {
				matches = append(matches, value)
			}
		}
	}

	if err := scanner.Err(); err != nil && len(matches) == 0 {
		values := grabDeclaredStringValues(input)
		for _, value := range values {
			match, _ := re.MatchString(value)
			if match {
				matches = append(matches, value)
			}
		}
	}

	matches = removeDuplicates(matches)

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
		//fmt.Printf("error reading string: %s\n", err)
	}

	return captured

}

func FindSecrets(text string) ToolData {

	var output ToolData
	var secrets []Secret

	var tags []string

	domains, _ := textsubs.DomainsOnly(text, false)
	domains = textsubs.Resolve(domains)

	for _, provider := range signatures {
		for service, regex := range provider.Keys {
			matches, err := getMatchingLines(text, regex)
			if err != nil {
				//log.Printf("Error reading data: %v\n", err)
				return output
			}

			if len(matches) > 0 {

				matches = removeDuplicates(matches)

				for _, match := range matches {

					tags = append(tags, "regexMatched")

					entropy := EntropyPercentage(match)
					if entropy > 66.6 {
						tags = append(tags, "highEntropy")
					}

					// todo: modify this to look at the variable named before the captured string and see if THAT
					// has the provider and/or service name or not.
					providerString := strings.ToLower(strings.Split(provider.Name, ".")[0])
					if strings.Contains(strings.ToLower(text), providerString) {
						tags = append(tags, "providerDetected")
					}

					secret := Secret{
						Provider:    provider.Name,
						ServiceName: service,
						Secret:      match,
						Entropy:     entropy,
						Tags:        removeDuplicates(tags),
					}

					// Remove this in the future when working with 
					// env files and stuff
					if strings.Contains(text, providerString) {
						secrets = append(secrets, secret)
					}

				}

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
			}
		}
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
