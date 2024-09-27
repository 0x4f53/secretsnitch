package main

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
	readInputFile("urls.txt")
	//readInputURL("https://meettaamaskrwwallet.gitbook.io/us")
}
