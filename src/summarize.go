package sdees

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func SummarizeEntries(texts []string, textsBranch []string) string {
	var summarized []string
	for i, text := range texts {
		dateInfo := strings.TrimSpace(strings.Split(text, " -==- ")[0])
		dateInfo = strings.Join(strings.Split(dateInfo, " ")[:5], " ")
		text = strings.Join(strings.Split(text, "\n")[1:], " ")
		words := strings.Split(text, " ")
		sentence := ""
		numWords := 1
		for {
			numWords++
			sentence = strings.Join(words[0:numWords], " ")
			if utf8.RuneCountInString(strings.TrimSpace(sentence)) > 80 || numWords >= len(words) {
				sentence = strings.Join(words[0:numWords-1], " ")
				break
			}
		}
		summarized = append(summarized,
				    fmt.Sprintf("%s - %s (%s words):\n  %s", HashIDToString(textsBranch[i]), dateInfo, NumberToString(len(words), ','), strings.TrimSpace(sentence)))
	}
	return strings.Join(summarized, "\n")
}

func NumberToString(n int, sep rune) string {

	s := strconv.Itoa(n)

	startOffset := 0
	var buff bytes.Buffer

	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}

	l := len(s)

	commaIndex := 3 - ((l - startOffset) % 3)

	if commaIndex == 3 {
		commaIndex = 0
	}

	for i := startOffset; i < l; i++ {

		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++

		buff.WriteByte(s[i])
	}

	return buff.String()
}
