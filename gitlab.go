package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://gitlab.com/api/v4"

type Namespace struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Namespace Namespace `json:"namespace"`
	WebURL    string    `json:"web_url"`
}

type Commit struct {
	ID         string    `json:"id"`
	ShortID    string    `json:"short_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorName string    `json:"author_name"`
	PatchURL   string    `json:"patch_url"`
}

type ProjectCommit struct {
	Project Project `json:"project"`
	Commit  Commit  `json:"commit"`
}

func getProjects(page int) ([]Project, error) {
	url := fmt.Sprintf("%s/projects?visibility=public&per_page=%d&page=%d", baseURL, perPage, page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch projects: %s", resp.Status)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func getCommits(projectID int) ([]Commit, error) {
	url := fmt.Sprintf("%s/projects/%d/repository/commits?per_page=%d", baseURL, projectID, perPage)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch commits for project %d: %s", projectID, resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	for i := range commits {
		commits[i].PatchURL = fmt.Sprintf("%s/projects/%d/repository/commits/%s.patch", baseURL, projectID, commits[i].ID)
	}
	return commits, nil
}

func getLast(perPage int, maxCommits int) []ProjectCommit {
	var allProjectCommits []ProjectCommit

	page := 1
	for len(allProjectCommits) < maxCommits {
		projects, err := getProjects(page)
		if err != nil {
			log.Fatalf("Error fetching projects: %v", err)
		}
		if len(projects) == 0 {
			break
		}

		for _, project := range projects {
			commits, err := getCommits(project.ID)
			if err != nil {
				log.Printf("Error fetching commits for project %d: %v", project.ID, err)
				continue
			}

			for _, commit := range commits {
				allProjectCommits = append(allProjectCommits, ProjectCommit{
					Project: project,
					Commit:  commit,
				})
				if len(allProjectCommits) >= maxCommits {
					break
				}
			}

			if len(allProjectCommits) >= maxCommits {
				break
			}
		}

		page++
	}

	return allProjectCommits

}
