package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/0x4f53/textsubs"
)

type Secret struct {
	Provider           string
	ServiceName        string
	Matches            []string
	CacheFile          string
	URL                string
	CapturedSubdomains []string
	//CapturedURLs    []string
	//Metadata        []string
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

func findSecrets(file string) Secret {

	var secret Secret

	subdomains, _ := textsubs.SubdomainsOnly(file, false)
	subdomains = textsubs.Resolve(subdomains)

	for _, service := range signatures {
		for keyName, regex := range service.Keys {
			matches, err := getMatchingLines(file, regex)
			if err != nil {
				//log.Printf("Error reading data: %v\n", err)
				return secret
			}

			if len(matches) > 0 {
				secret = Secret{
					Provider:           service.Name,
					ServiceName:        keyName,
					URL:                substringBeforeFirst(file, "---"),
					Matches:            matches,
					CapturedSubdomains: subdomains,
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

	secrets := findSecrets(string(data))
	secrets.CacheFile = filePath

	if len(secrets.Matches) > 0 {
		fmt.Println(secrets)
	}

}

func readCache(files []string) {

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
