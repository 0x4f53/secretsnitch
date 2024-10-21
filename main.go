package secretsnitch

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
		FetchFromUrlList(*urlList)
		files, _ := ListCachedFiles()
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

	if *github {
		patches := getGitHubPatchLinks(*to, *from)
		cacheGitHubPatchLinks(patches)
		FetchFromUrlList(gitHubPatchCache)
		files, _ := ListCachedFiles()
		ScanFiles(files)
		os.RemoveAll(gitHubPatchCache)
		os.RemoveAll(gitHubCommitsDirectory)
		return
	}

	if *gitlab {
		patches := getGitLabPatchLinks()
		cacheGitLabPatchLinks(patches)
		FetchFromUrlList(gitLabPatchCache)
		files, _ := ListCachedFiles()
		ScanFiles(files)
		os.RemoveAll(gitLabPatchCache)
		os.RemoveAll(gitLabCommitsDirectory)
		return
	}

	if *phishtank {
		savePhishtankDataset()
		FetchFromUrlList(phishtankURLCache)
		files, _ := ListCachedFiles()
		ScanFiles(files)
		os.RemoveAll(phishtankURLCache)
		return
	}

}
