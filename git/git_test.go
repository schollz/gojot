package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneral(t *testing.T) {
	gr, err := New("repo", "folder")
	assert.Equal(t, nil, err)
}
