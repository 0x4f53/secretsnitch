package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	maxWorkers = 100 // Number of concurrent workers
	semaphore  = make(chan struct{}, maxWorkers)
	wg         sync.WaitGroup // WaitGroup to wait for all workers to finish
)

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !info.IsDir() {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire a worker slot
			defer func() { <-semaphore }() // Release the worker slot when done

			dst := filepath.Join("temp", src)
			err := copyFile(src, dst)
			if err != nil {
				fmt.Printf("Error copying file %s: %v\n", src, err)
			} else {
				//fmt.Printf("Copied file %s to temp/%s\n", src, src)
			}
		}(path)
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create all directories in path to dst
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run copy_files.go <directory>")
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
