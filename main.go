package main

var (
	maxWorkers = 100000 // number of concurrent workers
)

var signatures []Signature

/*
func main() {
		if len(os.Args) < 2 {
			fmt.Println("Usage: go run main.go <directory>")
			return
		}
		root := os.Args[1]

		err := filepath.Walk(root, visit)
		if err != nil {
			fmt.Printf("Error walking the path %s: %v\n", root, err)
		}

		// Wait for all workers to finish
		wg.Wait()

}
*/

func main() {
	signatures = readSignatures()
	//readInputFile("urls.txt")
	files, _ := listCachedFiles()
	readCache(files)

}
