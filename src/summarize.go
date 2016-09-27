package gitsdees

import (
	"strings"
	"unicode/utf8"
)

func SummarizeEntries(cache Cache) string {
	var summarized []string
	texts, _ := CombineEntries(cache)
	for _, text := range texts {
		dateInfo := strings.TrimSpace(strings.Split(text, " -==- ")[0])
		dateInfo = strings.Join(strings.Split(dateInfo, " ")[:5], " ")
		text = strings.Join(strings.Split(text, "\n")[1:], " ")
		words := strings.Split(text, " ")
		sentence := ""
		numWords := 1
		for {
			numWords++
			sentence = strings.Join(words[0:numWords], " ")
			if utf8.RuneCountInString(dateInfo+" "+strings.TrimSpace(sentence)) > 80 || numWords >= len(words) {
				sentence = strings.Join(words[0:numWords-1], " ")
				break
			}
		}
		summarized = append(summarized, dateInfo+" "+strings.TrimSpace(sentence))
	}
	return strings.Join(summarized, "\n")
}
