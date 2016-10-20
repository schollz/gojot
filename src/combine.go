package sdees

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

func CombineEntries(cache Cache) ([]string, []string, map[string]string) {
	logger.Debug("Combining entries")
	var branchHashes = make(map[string]string)
	var data = make(map[string]combineData)
	var dateBranch = make(map[string]string) // map of dates to branch
	for branch := range cache.Branch {
		textData := HeadMatter(cache.Branch[branch].Date, branch) + cache.Branch[branch].Text
		branchHashes[branch] = GetMD5Hash(cache.Branch[branch].Text)
		parsedData, err := ParseDate(strings.TrimSpace(cache.Branch[branch].Date))
		if err != nil {
			logger.Debug(strings.TrimSpace(cache.Branch[branch].Date))
			logger.Error(err.Error())
		}
		data[branch] = combineData{date: parsedData, text: textData}
		dateBranch[parsedData.String()] = branch
	}

	sortedCombineData := make(timeSlice, 0, len(data))
	for _, d := range data {
		sortedCombineData = append(sortedCombineData, d)
	}
	sort.Sort(sortedCombineData)
	texts := make([]string, len(data))
	branches := make([]string, len(data))
	for i, val := range sortedCombineData {
		texts[i] = val.text
		branches[i] = dateBranch[val.date.String()]
	}
	return texts, branches, branchHashes
}
