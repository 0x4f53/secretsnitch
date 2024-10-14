package main

var maxWorkers = 100000 // number of concurrent workers

var signatures []Signature

func main() {

	setFlags()

	signatures = readSignatures()

	if *urlList != "" {
		fetchFromUrlList(*urlList)
		files, _ := listCachedFiles()
		ScanFiles(files)
	}

	//readInputFile("urls.txt")
	//files, _ := listCachedFiles()
	//readCache(files)

	// githubPatches.GetCommitsInRange("", from, to)

}
