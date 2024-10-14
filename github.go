package main

import (
	"github.com/0x4f53/github-patches"
)

var patchCache = ".githubPatchCache"
var commitsDirectory = ".githubCommits"

func getPatchLinks (to string, from string) []string {
	var patches []string
	githubPatches.GetCommitsInRange(commitsDirectory, to, from)
	files, _ := getAllFiles(commitsDirectory)
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

func cachePatchLinks(links []string) {
	for _, link := range links {
		appendToFile(patchCache, link)
	}
}