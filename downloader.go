/*
*
Worker-optimized downloading
Stress tested with 100k URLs from GitHub
*
*/

package main

import (
	"bufio"
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
		log.Println("Visiting", r.URL)
	})
	c.SetRequestTimeout(time.Duration(timeoutSeconds) * time.Second)
	c.OnResponse(func(r *colly.Response) {
		responseString := url + "\n---\n" + string(r.Body)

		err := saveToCache(makeCacheFilename(url), responseString)
		if err != nil {
			//log.Printf("Failed to write response body to file: %s\n", err)
		} else {
			//log.Printf("Content from %s saved to %s\n", url, cacheDir)
		}
	})

	err := c.Visit(url)
	if err != nil {
		//log.Printf("Failed to visit URL %s: %s\n", url, err)
	}

}

func fetchFromUrlListFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		//log.Fatalf("Failed to open file: %s\n", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)

	urlChan := make(chan string)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for url := range urlChan {
				cacheFileName := md5Hash(url)[0:8]
				if !fileExists(cacheDir + cacheFileName) {
					wg.Add(1)
					scrapeURL(url, &wg)
				}
			}
		}()
	}

	for scanner.Scan() {
		urlChan <- scanner.Text()
	}

	close(urlChan)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		//log.Fatalf("Error reading URLs: %s\n", err)
	}
}
func fetchFromUrlList(urls []string) {
	var wg sync.WaitGroup

	urlChan := make(chan string)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for url := range urlChan {
				if !fileExists(makeCacheFilename(url)) {
					wg.Add(1)
					scrapeURL(url, &wg)
				}
			}
		}()
	}

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
	wg.Wait()
}
