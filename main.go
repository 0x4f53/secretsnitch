package main

import (
	"os"
	"sync"
)

var maxWorkers = 100000 // number of concurrent workers

var signatures []Signature

func main() {

	setFlags()

	signatures = readSignatures()

	if *urlList != "" {
		fetchFromUrlList(*urlList)
		files, _ := listCachedFiles()
		ScanFiles(files)
		return
	}

	if *url != "" {
		var wg sync.WaitGroup
		cacheFileName := cacheDir + md5Hash(*url)[0:8] + cacheFileExtension
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

	/*
		if *github {
			patches := getGitHubPatchLinks(*to, *from)
			cacheGitHubPatchLinks(patches)
			fetchFromUrlList(gitHubPatchCache)
			files, _ := listCachedFiles()
			ScanFiles(files)
			os.RemoveAll(gitHubPatchCache)
			os.RemoveAll(gitHubCommitsDirectory)
			return
		}

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
		fetchFromUrlList(phishtankURLCache)
		files, _ := listCachedFiles()
		ScanFiles(files)
		os.RemoveAll(phishtankURLCache)
		return
	}

}
