package gogpg

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const secring = `C:\cygwin64\home\Zack\.gnupg\secring.gpg`
const pubring = `C:\cygwin64\home\Zack\.gnupg\pubring.gpg`

func BenchmarkDecryptFile(b *testing.B) {
	gs, _ := New(secring, pubring)
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
	gs, _ := New(secring, pubring)
	gs.Init("Testy McTestFace", "1234")
	enc, _ := gs.Encrypt([]byte("this is a test"))
	for n := 0; n < b.N; n++ {
		gs.Decrypt(enc)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	gs, _ := New(secring, pubring)
	gs.Init("Testy McTestFace", "1234")
	for n := 0; n < b.N; n++ {
		gs.Encrypt([]byte("this is a test"))
	}
}

func TestListing(t *testing.T) {
	gs, err := New(secring, pubring)
	assert.Equal(t, nil, err)
	keys, err := gs.ListPrivateKeys()
	assert.Equal(t, true, stringInSlice("Testy McTestFace", keys))
	assert.Equal(t, nil, err)

	keys, err = gs.ListPublicKeys()
	assert.Equal(t, true, stringInSlice("Testy McTestFace", keys))
	assert.Equal(t, nil, err)
}

func TestGeneral(t *testing.T) {
	gs, err := New(secring, pubring)
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
	gs, err := New(secring, pubring)
	assert.Equal(t, nil, err)
	err = gs.Init("Testy Blah", "1234")
	assert.Equal(t, NoSuchKeyError(NoSuchKeyError{key: "Testy Blah"}), err)
	err = gs.Init("Testy McTestFace", "12354")
	assert.Equal(t, IncorrectPassphrase(IncorrectPassphrase{key: "Testy McTestFace"}), err)
}
