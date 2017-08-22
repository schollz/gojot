package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneral(t *testing.T) {
	gr, err := New("https://github.com/schollz/schollz/asdf.git", "/tmp/testtest")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, nil != gr.Update())

	gr, err = New("https://github.com/schollz/schollz/test.git", "/tmp/testtest")
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, gr.Update())

	assert.Equal(t, "test", parseRepoFolder("git@github.com:schollz/test.git"))
	assert.Equal(t, "gojot", parseRepoFolder("https://github.com/schollz/gojot.git"))
}
