package main

import (
	"github.com/0x4f53/gitlab-patches"
)

var gitLabPatchCache = ".gitlabPatchCache"
var gitLabCommitsDirectory = ".gitlabCommits"

func getGitLabPatchLinks () []string {
	var patches []string
	patchData := gitlabPatches.GetGitlabCommits(100, 100)
	for _, patch := range patchData {
		patches = append(patches, patch.CommitPatchURL)
	}
	return patches
}

func cacheGitLabPatchLinks(links []string) {
	for _, link := range links {
		appendToFile(gitLabPatchCache, link)
	}
}