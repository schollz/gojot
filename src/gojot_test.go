package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGojotGeneral(t *testing.T) {
	os.RemoveAll(path.Join(cacheFolder, "demo2"))
	gj, err := New("https://github.com/schollz/demo2.git", true)
	assert.Nil(t, err)

	err = gj.Init()
	if err != nil {
		if err.Error() == "Need to make config file" {
			fmt.Println("Here's the available keys")
			fmt.Println(gj.gpg.ListPrivateKeys())
			id := "Testy McTestFace"
			passphrase := "1234"
			err2 := gj.gpg.Init(id, passphrase)
			assert.Nil(t, err2)
			err2 = gj.NewConfig()
			assert.Nil(t, err2)
		}
	}

	assert.Nil(t, gj.LoadConfig())
	assert.Equal(t, "Testy McTestFace", gj.config.Identity)
	assert.Equal(t, 4, strings.Count(gj.config.Salt, "-"))

	repos, err := ListDirs()
	assert.Nil(t, err)
	assert.Equal(t, true, strings.Contains(repos["https://github.com/schollz/demo2.git"], ".cache/gojot2/demo2"))
}
