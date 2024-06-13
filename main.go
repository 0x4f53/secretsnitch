package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	startDir := "."
	err := filepath.Walk(startDir, processFile)
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", startDir, err)
	}
}

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Printf("Contents of file %s:\n", path)
		fmt.Println(string(content))
	}
	return nil
}
