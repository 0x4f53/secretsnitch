package secretsnitch

import (
	"github.com/0x4f53/github-patches"
)

var gitHubPatchCache = ".githubPatchCache"
var gitHubCommitsDirectory = ".githubCommits"

func getGitHubPatchLinks (to string, from string) []string {
	var patches []string
	githubPatches.GetCommitsInRange(gitHubCommitsDirectory, to, from)
	files, _ := getAllFiles(gitHubCommitsDirectory)
	for _, file := range files {
		events, _ := githubPatches.ParseJSONFile(file)
		for _, event := range events {
			if event.PatchUrl != "" {
				patches = append(patches, event.PatchUrl)
			}
		} 
	}
	return patches
}

func cacheGitHubPatchLinks(links []string) {
	for _, link := range links {
		appendToFile(gitHubPatchCache, link)
	}
}