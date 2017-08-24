package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v1"
)

type Document struct {
	Front FrontMatter
	Text  string
}

type Documents []Document

func (p Documents) Len() int {
	return len(p)
}

func (p Documents) Less(i, j int) bool {
	return p[i].Front.LastModified.Time.Before(p[j].Front.LastModified.Time)
}

func (p Documents) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type YamlFrontMatter struct {
	Time            string
	LastModified    string `yaml:"last_modified" json:"last_modified"`
	Document, Entry string
	Tags            []string
}

type FrontMatter struct {
	Time            MyTime
	LastModified    MyTime `yaml:"last_modified" json:"last_modified"`
	Document, Entry string
	Tags            []string
}

type MyTime struct {
	time.Time
}

func (self *MyTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
	}
	self.Time = t
	return
}

func UnmarshalFrontMatter(b []byte) (fm FrontMatter, err error) {
	// Unmarshal the YAML
	yfm := YamlFrontMatter{}
	err = yaml.Unmarshal(b, &yfm)
	if err != nil {
		return
	}
	// YAML -> JSON
	jsonMarshalled, err := json.Marshal(yfm)
	if err != nil {
		return
	}
	// Unmarshal the JSON to correctly unmarshal time
	err = json.Unmarshal(jsonMarshalled, &fm)
	return
}

func MarshalFrontMatter(fm FrontMatter) (b []byte, err error) {
	// Build the YAML
	yfm := YamlFrontMatter{}
	yfm.Document = fm.Document
	yfm.Entry = fm.Entry
	yfm.LastModified = fm.LastModified.Format("2006-01-02 15:04:05")
	yfm.Time = fm.Time.Format("2006-01-02 15:04:05")
	yfm.Tags = fm.Tags
	// Marshal the YAML
	b, err = yaml.Marshal(yfm)
	return
}

func ParseScroll(fulltext string) (docs Documents, err error) {
	var doc Document

	numDocuments := 0
	for i, _ := range strings.Split(fulltext, "---") {
		if math.Mod(float64(i), 2) == float64(0) && i > 0 {
			numDocuments++
		}
	}

	docs = make(Documents, 0, numDocuments)
	for i, text := range strings.Split(fulltext, "---") {
		if math.Mod(float64(i), 2) == float64(0) && i > 0 {
			doc.Text = strings.TrimSpace(text)
			docs = append(docs, doc)
		} else if math.Mod(float64(i), 2) == float64(1) && i > 0 {
			doc.Front, err = UnmarshalFrontMatter([]byte(strings.TrimSpace(text)))
			if err != nil {
				fmt.Println(text)
				return
			}
		}
	}

	sort.Sort(docs)
	return
}
