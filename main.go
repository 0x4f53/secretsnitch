package main

var (
	maxWorkers = 100000 // number of concurrent workers
)

var signatures []Signature

func main() {
	//signatures = readSignatures()
	//readInputFile("urls.txt")
	files, _ := listCachedFiles()
	readCache(files)

}
