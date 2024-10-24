package main

import (
	"math"
)

// Function to calculate Shannon Entropy
func shannonEntropy(s string) float64 {
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}
	entropy := 0.0
	length := float64(len(s))
	for _, count := range freq {
		prob := count / length
		entropy -= prob * math.Log2(prob)
	}
	return entropy
}

// Function to calculate Rényi Entropy (alpha = 2)
func renyiEntropy(s string, alpha float64) float64 {
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}
	length := float64(len(s))
	sum := 0.0
	for _, count := range freq {
		prob := count / length
		sum += math.Pow(prob, alpha)
	}
	return (1 / (1 - alpha)) * math.Log2(sum)
}

// Function to calculate Min-Entropy
func minEntropy(s string) float64 {
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}
	maxProb := 0.0
	length := float64(len(s))
	for _, count := range freq {
		prob := count / length
		if prob > maxProb {
			maxProb = prob
		}
	}
	return -math.Log2(maxProb)
}

// Function to calculate Collision Entropy
func collisionEntropy(s string) float64 {
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}
	length := float64(len(s))
	sum := 0.0
	for _, count := range freq {
		prob := count / length
		sum += prob * prob
	}
	return -math.Log2(sum)
}

// Function to calculate Tsallis Entropy (q = 2)
func tsallisEntropy(s string, q float64) float64 {
	freq := make(map[rune]float64)
	for _, char := range s {
		freq[char]++
	}
	length := float64(len(s))
	sum := 0.0
	for _, count := range freq {
		prob := count / length
		sum += math.Pow(prob, q)
	}
	return (1 - sum) / (q - 1)
}

// Function to calculate entropy percentage from multiple entropy algorithms
func EntropyPercentage(text string) float64 {
	avgEntropy := AverageEntropy(text)
	uniqueChars := make(map[rune]bool)
	for _, char := range text {
		uniqueChars[char] = true
	}
	maxEntropy := math.Log2(float64(len(uniqueChars)))
	percentage := (avgEntropy / maxEntropy) * 100
	return percentage
}

// Function to calculate the average entropy score from multiple entropy algorithms
func AverageEntropy(text string) float64 {
	shannon := shannonEntropy(text)
	renyi := renyiEntropy(text, 2)
	min := minEntropy(text)
	collision := collisionEntropy(text)
	tsallis := tsallisEntropy(text, 2)

	average := (shannon + renyi + min + collision + tsallis) / 5
	return average
}
