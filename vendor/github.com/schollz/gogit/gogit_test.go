package gogit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommit(t *testing.T) {
	var err error
	gr, _ := New("https://github.com/schollz/test.git", "testtest")
	gr.Debug(true)
	err = gr.Update()
	assert.Equal(t, nil, err)
	err = gr.AddData([]byte("hi3"), "hello/hi3.txt")
	assert.Equal(t, nil, err)
	err = gr.Push()
	assert.Equal(t, nil, err)
}

func TestClone(t *testing.T) {
	gr, err := New("https://github.com/schollz/test.git", "testtest")
	gr.Debug(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, gr.Update())
	assert.Equal(t, true, exists("testtest"))
	assert.Equal(t, nil, gr.Update())
	os.RemoveAll("testtest")
}
func TestGeneral(t *testing.T) {
	gr, err := New("git@github.com:schollz/asdf.git", "testtest")
	gr.Debug(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, nil != gr.Update())

	gr, err = New("git@github.com:schollz/test.git")
	gr.Debug(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, gr.Update())
	assert.Equal(t, true, exists("test"))
	assert.Equal(t, nil, gr.Update())

	origin, err := GetRemoteOriginURL("test")
	assert.Equal(t, nil, err)
	assert.Equal(t, "git@github.com:schollz/test.git", origin)
	os.RemoveAll("testtest")
	os.RemoveAll("test")

	assert.Equal(t, "test", ParseRepoFolder("git@github.com:schollz/test.git"))
	assert.Equal(t, "gojot", ParseRepoFolder("https://github.com/schollz/gojot.git"))
}
