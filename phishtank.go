package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
)

var phishtankURLCache = "./phishtankURLCache"
var phishtankURL = "http://data.phishtank.com/data/online-valid.csv"

func savePhishtankDataset() error {
	resp, err := http.Get(phishtankURL)

	switch resp.StatusCode {
	case 404:
		fmt.Println("Phishtank is down. Please try again later.")
		os.Exit(-1)
	case 429:
		fmt.Println("Phishtank has rate-limited you. Please try again.")
		os.Exit(-1)
	}

	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	reader := csv.NewReader(resp.Body)

	outputFile, err := os.Create(phishtankURLCache)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record: %w", err)
		}

		if len(record) >= 2 {
			if _, err := outputFile.WriteString(record[1] + "\n"); err != nil {
				return fmt.Errorf("failed to write to output file: %w", err)
			}
		}
	}

	return nil
}
