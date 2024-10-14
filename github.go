package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var githubCacheDir = ".githubCache/"

func printTimestamps(from, to string) []string {

	var timestamps []string

	if from == "" && to == "" {
		now := time.Now()
		previousHour := now.Add(-1 * time.Hour)
		timestamps = append(timestamps, previousHour.Format("2006-01-02-3"))
		return timestamps
	}

	fromTime, err := time.Parse("2006-01-02-3", from)
	if err != nil {
		fmt.Println("Invalid 'from' timestamp format. Use dd-mm-yyyy-H.")
		return timestamps
	}

	toTime, err := time.Parse("2006-01-02-3", to)
	if err != nil {
		fmt.Println("Invalid 'to' timestamp format. Use dd-mm-yyyy-H.")
		return timestamps
	}

	var step time.Duration
	if fromTime.After(toTime) {
		step = -1 * time.Hour
	} else {
		step = 1 * time.Hour
	}

	for t := fromTime; ; t = t.Add(step) {
		timestamps = append(timestamps, t.Format("2006-01-02-3"))
		if (step > 0 && !t.Before(toTime)) || (step < 0 && !t.After(toTime)) {
			break
		}
	}

	return timestamps
}

func printGharchiveChunkUrls(from, to string) []string {
	timestamps := printTimestamps(from, to)
	var urls []string

	for _, timestamp := range timestamps {
		url := fmt.Sprintf("https://data.gharchive.org/%s.json.gz", timestamp)
		urls = append(urls, url)
	}

	return urls

}

func downloadAndExtract(url string) error {

	makeDir(githubCacheDir)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	fileName := filepath.Base(url)
	outFile, err := os.Create(githubCacheDir + fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if err = outFile.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	gzFileName := githubCacheDir + fileName

	gzFile, err := os.Open(gzFileName)

	if err != nil {
		return fmt.Errorf("failed to open gz file: %w", err)
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	jsonFileName := gzFileName[:len(gzFileName)-3]
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return fmt.Errorf("failed to create json file: %w", err)
	}
	defer jsonFile.Close()

	if _, err = io.Copy(jsonFile, gzReader); err != nil {
		return fmt.Errorf("failed to write json file: %w", err)
	}

	fmt.Printf("Extracted %s to %s\n", gzFileName, jsonFileName)
	os.Remove(gzFileName)

	return nil

}

func makeDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create directory: %w", err)
		}
	}
	return nil
}

// Root represents the overall structure of the JSON
type Root struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Actor     Actor        `json:"actor"`
	Repo      Repo         `json:"repo"`
	Payload   Payload      `json:"payload"`
	Public    bool         `json:"public"`
	CreatedAt string       `json:"created_at"`
	Org       Organization `json:"org,omitempty"` // Optional field
}

// Actor represents the user who triggered the event
type Actor struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	GravatarID string `json:"gravatar_id"`
	URL        string `json:"url"`
	AvatarURL  string `json:"avatar_url"`
}

// Repo represents the repository involved in the event
type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Payload represents the data specific to the event
type Payload struct {
	Action       string   `json:"action,omitempty"`        // Present in IssuesEvent
	Issue        *Issue   `json:"issue,omitempty"`         // Present in IssuesEvent
	Ref          string   `json:"ref,omitempty"`           // Present in CreateEvent
	RefType      string   `json:"ref_type,omitempty"`      // Present in CreateEvent
	MasterBranch string   `json:"master_branch,omitempty"` // Present in CreateEvent
	Commits      []Commit `json:"commits,omitempty"`       // Present in PushEvent
}

// Issue represents an issue in the repository
type Issue struct {
	URL         string     `json:"url"`
	LabelsURL   string     `json:"labels_url"`
	CommentsURL string     `json:"comments_url"`
	EventsURL   string     `json:"events_url"`
	HTMLURL     string     `json:"html_url"`
	ID          int        `json:"id"`
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	User        Actor      `json:"user"`
	Labels      []Label    `json:"labels"`
	State       string     `json:"state"`
	Locked      bool       `json:"locked"`
	Assignee    *Actor     `json:"assignee,omitempty"`
	Milestone   *Milestone `json:"milestone,omitempty"`
	Comments    int        `json:"comments"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	ClosedAt    *string    `json:"closed_at,omitempty"`
	Body        string     `json:"body"`
}

// Label represents a label for an issue
type Label struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Milestone represents a milestone in the repository
type Milestone struct {
	URL          string  `json:"url"`
	HTMLURL      string  `json:"html_url"`
	ID           int     `json:"id"`
	Number       int     `json:"number"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Creator      Actor   `json:"creator"`
	OpenIssues   int     `json:"open_issues"`
	ClosedIssues int     `json:"closed_issues"`
	State        string  `json:"state"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DueOn        *string `json:"due_on,omitempty"`
	ClosedAt     *string `json:"closed_at,omitempty"`
}

// Commit represents a commit in the PushEvent payload
type Commit struct {
	SHA     string `json:"sha"`
	Author  Author `json:"author"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

// Author represents the author of a commit
type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Organization represents the organization related to the event
type Organization struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	GravatarID string `json:"gravatar_id"`
	URL        string `json:"url"`
	AvatarURL  string `json:"avatar_url"`
}

func getCommitUrls() {
	jsonData := `[{"id":"3487637759","type":"IssuesEvent",...}]` // Your full JSON data here

	var events []Root
	err := json.Unmarshal([]byte(jsonData), &events)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Print the parsed data
	for _, event := range events {
		fmt.Printf("Event ID: %s, Type: %s, Actor: %s\n", event.ID, event.Type, event.Actor.Login)
	}
}

func main() {
	urls := printGharchiveChunkUrls("2016-01-02-7", "2016-01-02-3")
	fmt.Println(urls)

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := downloadAndExtract(url); err != nil {
				fmt.Printf("Error processing %s: %v\n", url, err)
			}
		}(url)
	}
	wg.Wait()

}
