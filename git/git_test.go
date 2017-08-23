package git

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneral(t *testing.T) {
	gr, err := New("https://github.com/schollz/asdf.git", "testtest")
	gr.Debug(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, nil != gr.Update())

	gr, err = New("https://github.com/schollz/test.git", "testtest")
	gr.Debug(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, gr.Update())
	assert.Equal(t, true, exists("testtest"))
	assert.Equal(t, true, exists("testtest/test"))
	assert.Equal(t, nil, gr.Update())

	origin, err := GetRemoteOriginURL("testtest/test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://github.com/schollz/test.git", origin)
	os.RemoveAll("testtest")

	assert.Equal(t, "test", parseRepoFolder("git@github.com:schollz/test.git"))
	assert.Equal(t, "gojot", parseRepoFolder("https://github.com/schollz/gojot.git"))
}
