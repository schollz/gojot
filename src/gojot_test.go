package main

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocuments(t *testing.T) {
	// This test uses one of my old entries to verify that there is compatibility
	doc := `---

time: '2017-02-14 16:01:39'
last_modified: '2017-06-09 20:00:36'
document: ideas
entry: '44982409'

---
Are proteins ever in equilibrium? Is it possible that they only need to exist
"natively" for the lifetime of the human or the human cell?`
	os.RemoveAll(path.Join(cacheFolder, "demo2"))
	gj, err := New(true)
	gj.SetRepo("https://github.com/schollz/demo2.git")
	gj.config.Salt = "39c7c512-25be-4a35-921c-629f5d67fd88"
	assert.Nil(t, err)
	docs, err := gj.ParseDocuments(doc)
	assert.Nil(t, err)
	assert.Equal(t, "8ATeHRiksk/70c8d449d188d0824f95da8eb693c14c.asc", docs[0].file)
}

func TestGojotGeneral(t *testing.T) {
	os.RemoveAll(path.Join(cacheFolder, "demo2"))
	gj, err := New(true)
	assert.Nil(t, err)

	err = gj.SetRepo("https://github.com/schollz/demo2.git")
	id := "Testy McTestFace"
	passphrase := "1234"
	err = gj.LoadConfig(id, passphrase)
	assert.Nil(t, err)

	assert.Equal(t, "Testy McTestFace", gj.config.Identity)
	assert.Equal(t, 4, strings.Count(gj.config.Salt, "-"))

	repos, err := ListAvailableRepos()
	assert.Nil(t, err)
	assert.Equal(t, true, strings.Contains(repos["https://github.com/schollz/demo2.git"], ".cache/gojot2/demo2"))
}
