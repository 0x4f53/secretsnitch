package main

import (
	"fmt"
	"os"
	"sync"

	githubPatches "github.com/0x4f53/github-patches"
)

var maxWorkers = 100000 // number of concurrent workers

var signatures []Signature

func main() {

	setFlags()

	signatures = readSignatures()

	if *urlList != "" {
		fetchFromUrlListFile(*urlList)
		files, _ := listCachedFiles()
		ScanFiles(files)
		return
	}

	if *url != "" {

		cacheFileName := makeCacheFilename(*url)

		var wg sync.WaitGroup

		if !fileExists(cacheFileName) {
			wg.Add(1)
			scrapeURL(*url, &wg)
		}

		wg.Add(1)
		scanFile(cacheFileName, &wg)
		wg.Wait()
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
		patchFiles, _ := listFiles(githubPatches.GithubCacheDir)
		parsedData, _ := githubPatches.ParseJSONFiles(patchFiles)

		fmt.Println(parsedData)

		var patches []string
		for _, patch := range parsedData {
			patches = append(patches, patch.PatchUrl)
		}

		fetchFromUrlList(patches)

		files, _ := listCachedFiles()
		ScanFiles(files)
		os.RemoveAll(githubPatches.GithubCacheDir)
		return
	}

	/*
		if *gitlab {
			patches := getGitLabPatchLinks()
			cacheGitLabPatchLinks(patches)
			fetchFromUrlList(gitLabPatchCache)
			files, _ := listCachedFiles()
			ScanFiles(files)
			os.RemoveAll(gitLabPatchCache)
			os.RemoveAll(gitLabCommitsDirectory)
			return
		}
	*/

	if *phishtank {
		savePhishtankDataset()
		fetchFromUrlListFile(phishtankURLCache)
		files, _ := listCachedFiles()
		ScanFiles(files)
		os.RemoveAll(phishtankURLCache)
		return
	}

}
