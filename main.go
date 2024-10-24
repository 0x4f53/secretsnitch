package main

import (
	"os"

	githubPatches "github.com/0x4f53/github-patches"
	gitlabPatches "github.com/0x4f53/gitlab-patches"
)

var maxWorkers = 100000 // number of concurrent workers

var signatures []Signature

func main() {

	setFlags()

	signatures = readSignatures()

	if *urlList != "" {
		urls, _ := readLines(*urlList)
		successfulUrls := fetchFromUrlList(urls)
		ScanFiles(successfulUrls)
		return
	}

	if *url != "" {
		successfulUrls := fetchFromUrlList([]string{*url})
		ScanFiles(successfulUrls)
		return
	}

	if *directory != "" {
		files, _ := getAllFiles(*directory)
		ScanFiles(files)
		return
	}

	if *file != "" {
		ScanFiles([]string{*file})
		return
	}

	if *github {
		githubPatches.GetCommitsInRange(githubPatches.GithubCacheDir, *from, *to, false)
		chunks, _ := listFiles(githubPatches.GithubCacheDir)

		var patches []string

		for _, chunk := range chunks {
			events, _ := githubPatches.ParseGitHubCommits(githubPatches.GithubCacheDir + chunk)

			for _, event := range events {
				for _, commit := range event.Payload.Commits {
					patches = append(patches, commit.PatchURL)
				}
			}

		}

		successfulUrls := fetchFromUrlList(patches)
		ScanFiles(successfulUrls)
		os.RemoveAll(githubPatches.GithubCacheDir)
		return
	}

	if *gitlab {
		commitData := gitlabPatches.GetGitlabCommits(100, 100)

		var patches []string
		for _, patch := range commitData {
			patches = append(patches, patch.CommitPatchURL)
		}

		successfulUrls := fetchFromUrlList(patches)
		ScanFiles(successfulUrls)
		os.RemoveAll(gitlabPatches.GitlabCacheDir)
		return
	}

	if *githubGists {
		gistData := githubPatches.GetLast100Gists()
		parsedGists, _ := githubPatches.ParseGistData(gistData)

		var gists []string
		for _, gist := range parsedGists {
			gists = append(gists, gist.RawURL)
		}

		successfulUrls := fetchFromUrlList(gists)
		ScanFiles(successfulUrls)
		return
	}

	if *phishtank {
		savePhishtankDataset()
		urls, _ := readLines(phishtankURLCache)
		successfulUrls := fetchFromUrlList(urls)
		ScanFiles(successfulUrls)
		return
	}

}
