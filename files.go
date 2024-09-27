/*
*
Worker-optimized downloading
Stress tested with 100k URLs from GitHub - took around 2 minutes on an i5-8350U / 16GB DDR3 @ 1333 MHz /
*
*/

package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
)

var (
	cacheDir       = ".urlCache/"
	cacheExtension = ".cache"
)

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		lines = append(lines, host)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func md5Hash(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func makeDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create directory: %w", err)
		}
	}
	return nil
}

func cacheSave(filename string, data string) error {
	// Check if cache directory exists
	makeDir(cacheDir)
	filename = cacheDir + filename + cacheExtension
	err := os.WriteFile(filename, []byte(data), 0644)
	return err
}

func fileExists(location string) bool {
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		return true
	}
	return false
}