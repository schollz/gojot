package gogpg

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDecryptFile(b *testing.B) {
	gs, _ := New(false)
	gs.Init("Testy McTestFace", "1234")
	for n := 0; n < b.N; n++ {
		data, err := ioutil.ReadFile("testing/hello.txt.asc")
		if err != nil {
			panic(err)
		}
		gs.Decrypt(data)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	gs, _ := New(false)
	gs.Init("Testy McTestFace", "1234")
	enc, _ := gs.Encrypt([]byte("this is a test"))
	for n := 0; n < b.N; n++ {
		gs.Decrypt(enc)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	gs, _ := New(false)
	gs.Init("Testy McTestFace", "1234")
	for n := 0; n < b.N; n++ {
		gs.Encrypt([]byte("this is a test"))
	}
}

func TestListing(t *testing.T) {
	gs, err := New(true)
	gs.Debug(true)
	assert.Equal(t, nil, err)
	keys, err := gs.ListPrivateKeys()
	assert.Equal(t, true, stringInSlice("Testy McTestFace", keys))
	assert.Equal(t, nil, err)

	keys, err = gs.ListPublicKeys()
	assert.Equal(t, true, stringInSlice("Testy McTestFace", keys))
	assert.Equal(t, nil, err)
}

func TestGeneral(t *testing.T) {
	gs, err := New(true)
	assert.Equal(t, nil, err)
	err = gs.Init("Testy McTestFace", "1234")
	assert.Equal(t, nil, err)
	data, _ := ioutil.ReadFile("testing/hello.txt.asc")
	decrypted, err := gs.Decrypt(data)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Hello, world.\n", string(decrypted))

	encrypted, err := gs.Encrypt([]byte("Hello, world.\n"))
	assert.Equal(t, nil, err)
	decrypted, err = gs.Decrypt(encrypted)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Hello, world.\n", string(decrypted))
}

func TestBugs(t *testing.T) {
	gs, err := New(true)
	assert.Equal(t, nil, err)
	err = gs.Init("Testy Blah", "1234")
	assert.Equal(t, NoSuchKeyError(NoSuchKeyError{key: "Testy Blah"}), err)
	err = gs.Init("Testy McTestFace", "12354")
	assert.Equal(t, IncorrectPassphrase(IncorrectPassphrase{key: "Testy McTestFace"}), err)
}

func TestBulk(t *testing.T) {
	gs, _ := New(true)
	gs.Init("Testy McTestFace", "1234")
	filenames := make([]string, 1000)
	for i := 0; i < len(filenames); i++ {
		filenames[i] = "bulk/file" + strconv.Itoa(i) + ".asc"
	}
	if !exists("bulk") {
		os.MkdirAll("bulk", 0700)
	}
	for _, filename := range filenames {
		enc, err := gs.Encrypt([]byte(RandStringBytes(1000)))
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filename, enc, 0644)
		if err != nil {
			panic(err)
		}
	}
	data, err := gs.BulkDecrypt(filenames)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1000, len(data))
	os.RemoveAll("bulk")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
