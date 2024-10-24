/*
*
Worker-optimized downloading
Stress tested with 100k URLs from GitHub
*
*/

package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var (
	timeoutSeconds  = 10
	userAgentString = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
)

func scrapeURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgentString)
	})
	c.SetRequestTimeout(time.Duration(timeoutSeconds) * time.Second)
	c.OnResponse(func(r *colly.Response) {
		responseString := url + "\n---\n" + string(r.Body)

		makeDir(cacheDir)
		err := os.WriteFile(makeCacheFilename(url), []byte(responseString), 0644)

		if err != nil {
			log.Printf("Failed to write response body to file: %s\n", err)
		} else {
			log.Printf("Content from %s saved to %s\n", url, cacheDir)
		}
	})

	err := c.Visit(url)
	if err != nil {
		//log.Printf("Failed to visit URL %s: %s\n", url, err)
	}

}

func fetchFromUrlList(urls []string) []string {
	var wg sync.WaitGroup

	urlChan := make(chan string)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for url := range urlChan {
				if fileExists(makeCacheFilename(url)) {
					log.Printf("Skipping %s as it is already cached", url)
					continue
				}
				wg.Add(1)
				log.Printf("Visiting %s", url)
				scrapeURL(url, &wg)
			}
		}()
	}

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
	wg.Wait()

	var successfulDownloads []string

	cachedFiles, _ := listCachedFiles()
	for _, url := range urls {
		cachedFileName := makeCacheFilename(url)
		if sliceContainsString(cachedFiles, cachedFileName) {
			successfulDownloads = append(successfulDownloads, cachedFileName)
		}
	}

	return successfulDownloads

}
