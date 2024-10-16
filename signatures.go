package main

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var (
	signaturesFile = "signatures.yaml"
)

type Signature struct {
	Name string
	Keys map[string]string
}

func readSignatures() []Signature {

	var services []map[string][]map[string]string

	data, err := os.ReadFile(signaturesFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var parsedSignatures []Signature

	for _, serviceMap := range services {
		for serviceName, keys := range serviceMap {
			keyMap := make(map[string]string)
			for _, keyPair := range keys {
				for keyName, keyValue := range keyPair {
					keyMap[keyName] = keyValue
				}
			}
			parsedSignatures = append(parsedSignatures, Signature{
				Name: serviceName,
				Keys: keyMap,
			})
		}
	}

	return parsedSignatures

}

func expressionValues(inputMap map[string]*regexp.Regexp) []*regexp.Regexp {
	var values []*regexp.Regexp
	for _, value := range inputMap {
		values = append(values, value)
	}
	return values
}

func expressionKeys(inputMap map[string]*regexp.Regexp) []string {
	var keys []string
	for key, _ := range inputMap {
		keys = append(keys, key)
	}
	return keys
}

func expressionValue(inputMap map[string]*regexp.Regexp) *regexp.Regexp {
	return expressionValues(inputMap)[0]
}

func expressionKey(inputMap map[string]*regexp.Regexp) string {
	return expressionKeys(inputMap)[0]
}
