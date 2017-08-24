package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	salt := "39c7c512-25be-4a35-921c-629f5d67fd88"
	enc, err := Encode("notes", salt)
	assert.Equal(t, nil, err)
	assert.Equal(t, "BWSxs4f6iB", enc)

	dec, err := Decode(enc, salt)
	assert.Equal(t, nil, err)
	assert.Equal(t, "notes", dec)
}
