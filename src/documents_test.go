package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScroll(t *testing.T) {
	fulltext := `---	
time: '2016-02-17 06:34:59'
last_modified: '2017-06-09 20:00:38'
document: doco1
entry: entro1
---

This is some text 

---	
time: '2015-02-16 06:34:59'
last_modified: '2015-02-16 06:34:59'
document: doco1
entry: entro0
---

First entry

---	
time: '2016-02-17 06:34:59'
last_modified: '2016-02-17 06:34:59'
document: doco1
entry: entro1
---

This is some


`
	docs, err := ParseScroll(fulltext)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(docs))

	// Test whether Documents can be marhsalled/unmarshalled
	b, _ := json.Marshal(docs)
	var docs2 Documents
	json.Unmarshal(b, &docs2)
	assert.Equal(t, docs, docs)

	docString, err := docs[0].String()
	assert.Nil(t, err)
	assert.Equal(t, "---\ntime: 2015-02-16 06:34:59\nlast_modified: 2015-02-16 06:34:59\ndocument: doco1\nentry: entro0\ntags: []\n---\n\nFirst entry", docString)

	docsString, err := docs.String()
	assert.Nil(t, err)
	assert.Equal(t, "---\ntime: 2015-02-16 06:34:59\nlast_modified: 2015-02-16 06:34:59\ndocument: doco1\nentry: entro0\ntags: []\n---\n\nFirst entry\n\n---\ntime: 2016-02-17 06:34:59\nlast_modified: 2016-02-17 06:34:59\ndocument: doco1\nentry: entro1\ntags: []\n---\n\nThis is some\n\n---\ntime: 2016-02-17 06:34:59\nlast_modified: 2017-06-09 20:00:38\ndocument: doco1\nentry: entro1\ntags: []\n---\n\nThis is some text", docsString)

	docsString, err = docs.String("doco1")
	fmt.Println("+++++++")
	fmt.Println(docsString)
	fmt.Println("+++++++")
	assert.Nil(t, err)
	assert.Equal(t, "---\ntime: 2015-02-16 06:34:59\nlast_modified: 2015-02-16 06:34:59\ndocument: doco1\nentry: entro0\ntags: []\n---\n\nFirst entry\n\n---\ntime: 2016-02-17 06:34:59\nlast_modified: 2017-06-09 20:00:38\ndocument: doco1\nentry: entro1\ntags: []\n---\n\nThis is some text", docsString)
}

func TestScrollFrontMatter(t *testing.T) {
	header := `time: '2017-02-17 06:34:59'
last_modified: '2017-06-09 20:00:38'
document: doco1
entry: entro1
`
	h, err := UnmarshalFrontMatter([]byte(header))
	assert.Equal(t, nil, err)
	assert.Equal(t, "2017-02-17 06:34:59 +0000 UTC", h.Time.String())
	assert.Equal(t, "2017-06-09 20:00:38 +0000 UTC", h.LastModified.String())
	assert.Equal(t, "doco1", h.Document)
	assert.Equal(t, "entro1", h.Entry)

	marshalled, err := MarshalFrontMatter(h)
	assert.Equal(t, nil, err)
	assert.Equal(t, "time: 2017-02-17 06:34:59\nlast_modified: 2017-06-09 20:00:38\ndocument: doco1\nentry: entro1\ntags: []\n", string(marshalled))

	h, err = UnmarshalFrontMatter(marshalled)
	assert.Equal(t, nil, err)
	assert.Equal(t, "2017-02-17 06:34:59 +0000 UTC", h.Time.String())
	assert.Equal(t, "2017-06-09 20:00:38 +0000 UTC", h.LastModified.String())
	assert.Equal(t, "doco1", h.Document)
	assert.Equal(t, "entro1", h.Entry)

}
