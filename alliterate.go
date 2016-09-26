package main

import (
	"math/rand"
	"strings"
	"time"
)

func MakeAlliteration() string {
	dataAdj, _ := Asset("bin/adjectives.txt")
	dataNoun, _ := Asset("bin/nouns.txt")
	adjectives := strings.Split(string(dataAdj), "\n")
	nouns := strings.Split(string(dataNoun), "\n")
	alliterate := make(map[string]map[string][]string)
	for _, word := range adjectives {
		word = strings.Title(strings.TrimSpace(word))
		char0 := word[0:1] + " "
		if _, ok := alliterate[char0]; !ok {
			alliterate[char0] = make(map[string][]string)
		}
		if _, ok := alliterate[char0]["adjectives"]; !ok {
			alliterate[char0]["adjectives"] = []string{}
		}
		alliterate[char0]["adjectives"] = append(alliterate[char0]["adjectives"], word)
	}

	for _, word := range nouns {
		word = strings.Title(strings.TrimSpace(word))
		char0 := word[0:1] + " "
		if _, ok := alliterate[char0]; !ok {
			continue
		}
		if _, ok := alliterate[char0]["nouns"]; !ok {
			alliterate[char0]["nouns"] = []string{}
		}
		alliterate[char0]["nouns"] = append(alliterate[char0]["nouns"], word)
	}

	for _, word := range adjectives {
		word = strings.Title(strings.TrimSpace(word))
		char0 := word[0:1] + " "
		if _, ok := alliterate[char0]["nouns"]; !ok {
			delete(alliterate, char0)
		}
	}

	rand.Seed(time.Now().UnixNano())
	randomLetterID := rand.Intn(len(alliterate))
	i := 0
	randomLetter := "a "
	for val := range alliterate {
		if i == randomLetterID {
			randomLetter = val
		}
		i++
	}
	nounLength := len(alliterate[randomLetter]["nouns"])
	adjectiveLength := len(alliterate[randomLetter]["adjectives"])
	return alliterate[randomLetter]["adjectives"][rand.Intn(adjectiveLength)] + alliterate[randomLetter]["nouns"][rand.Intn(nounLength)]
}
