package main

import (
	"encoding/json"
	"math"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v1"
)

type Document struct {
	Front FrontMatter
	Text  string
	hash  string
	file  string
}

func NewDocument(document, entry string) (d *Document) {
	d = new(Document)
	d.Text = ""
	d.Front = FrontMatter{
		Document: document,
		Entry:    entry,
		Time: MyTime{
			time.Now(),
		},
		LastModified: MyTime{
			time.Now(),
		},
		Tags: []string{},
	}
	return
}

func (d *Document) String() (s string, err error) {
	s = "---\n"
	fm, err := MarshalFrontMatter(d.Front)
	if err != nil {
		return "", err
	}
	s += string(fm)
	s += "---\n\n"
	s += d.Text
	return
}

type Documents []Document

func (d Documents) String(documentFilter ...string) (s string, err error) {
	m := make(map[string]bool)
	docStrings := make([]string, d.Len())
	docStringI := 0
	for i := d.Len() - 1; i >= 0; i-- {
		if len(documentFilter) > 0 {
			// only add if in document
			if d[i].Front.Document != documentFilter[0] {
				continue
			}
			// only add if it is new
			if _, ok := m[d[i].Front.Entry]; ok {
				continue
			}
		}
		m[d[i].Front.Entry] = true
		docStrings[docStringI], err = d[i].String()
		if err != nil {
			return "", err
		}
		docStringI++
	}
	docStrings = docStrings[:docStringI]
	for i, j := 0, len(docStrings)-1; i < j; i, j = i+1, j-1 {
		docStrings[i], docStrings[j] = docStrings[j], docStrings[i]
	}
	s = strings.Join(docStrings, "\n\n")
	return
}

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
				return
			}
		}
	}

	sort.Sort(docs)
	return
}
