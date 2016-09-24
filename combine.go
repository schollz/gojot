package main

import (
	"sort"
	"strings"
	"time"
)

type combineData struct {
	date time.Time
	text string
}

type timeSlice []combineData

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return p[i].date.Before(p[j].date)
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func CombineEntries(cache map[string]Entry) []string {
	logger.Debug("Combining entries")
	var data = make(map[string]combineData)
	for branch := range cache {
		if cache[branch].Document == CurrentDocument {
			textData := HeadMatter(cache[branch].Date, cache[branch].Branch, cache[branch].Text) + cache[branch].Text
			parsedData, err := ParseDate(strings.TrimSpace(cache[branch].Date))
			if err != nil {
				logger.Debug(strings.TrimSpace(cache[branch].Date))
				logger.Error(err.Error())
			}
			data[branch] = combineData{date: parsedData, text: textData}
		}
	}

	sortedCombineData := make(timeSlice, 0, len(data))
	for _, d := range data {
		sortedCombineData = append(sortedCombineData, d)
	}
	sort.Sort(sortedCombineData)
	texts := make([]string, len(data))
	for i, val := range sortedCombineData {
		texts[i] = val.text
	}
	return texts
}
