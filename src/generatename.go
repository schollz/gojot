package sdees

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateEntryName() string {
	return GenerateNonAlliteration()
}

func GenerateAlliteration() string {
	dataAdj, _ := Asset("bin/adjectives.txt")
	dataNoun, _ := Asset("bin/nouns.txt")
	adjectives := strings.Split(string(dataAdj), "\n")
	nouns := strings.Split(string(dataNoun), "\n")
	adjectives = adjectives[0 : len(adjectives)-1]
	nouns = nouns[0 : len(nouns)-1]
	alliterate := make(map[string]map[string][]string)
	for _, word := range adjectives {
		word = strings.Title(strings.TrimSpace(word))
		if len(word) < 2 {
			continue
		}
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
		if len(word) < 2 {
			continue
		}
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
		if len(word) < 2 {
			continue
		}
		char0 := word[0:1] + " "
		if _, ok := alliterate[char0]["nouns"]; !ok {
			delete(alliterate[char0], "adjectives")
			delete(alliterate, char0)
		}
	}

	// // Count how many are possible
	// possible := 0
	// for letter := range alliterate {
	// 	possible += len(alliterate[letter]["nouns"]) * len(alliterate[letter]["adjectives"])
	// }
	// fmt.Printf("There are %d possible alliterative combinations\n", possible)
	// fmt.Println("There are %d total possible combinations\n", len(adjectives)*len(adjectives)*len(nouns))

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
	return alliterate[randomLetter]["adjectives"][rand.Intn(adjectiveLength)] + alliterate[randomLetter]["nouns"][rand.Intn(nounLength)] + strconv.Itoa(rand.Intn(100))
}

func GenerateNonAlliteration() string {
	dataAdj, _ := Asset("bin/adjectives.txt")
	dataNoun, _ := Asset("bin/nouns.txt")
	adjectives0 := strings.Split(string(dataAdj), "\n")
	nouns0 := strings.Split(string(dataNoun), "\n")
	adjectives := make([]string, len(adjectives0))
	nouns := make([]string, len(nouns0))
	j := 0
	for _, adjective := range adjectives0 {
		adjective = strings.TrimSpace(adjective)
		if len(adjective) > 3 {
			adjectives[j] = strings.Title(adjective)
			j++
		}
	}
	adjectives = adjectives[0:j]
	j = 0
	for _, noun := range nouns0 {
		noun = strings.TrimSpace(noun)
		if len(noun) > 3 {
			nouns[j] = strings.Title(noun)
			j++
		}
	}
	nouns = nouns[0:j]
	// fmt.Println("There are %d total possible combinations\n", len(adjectives)*len(adjectives)*len(nouns))

	rand.Seed(time.Now().UnixNano())
	randomNoun := nouns[rand.Intn(len(nouns))]
	randomAdjective1 := adjectives[rand.Intn(len(adjectives))]
	randomAdjective2 := adjectives[rand.Intn(len(adjectives))]
	// fmt.Printf("\n[%s],[%s],[%s]\n", randomNoun, randomAdjective1, randomAdjective2)
	return randomAdjective1 + randomAdjective2 + randomNoun
}
