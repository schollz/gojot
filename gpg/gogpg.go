package gogpg

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// basic functions
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GPGStore is the basic store object.
type GPGStore struct {
	secretKeyring     string
	publicKeyring     string
	identity          string
	passphrase        string
	privateEntity     *openpgp.Entity
	privateKeyToUse   []*openpgp.Entity
	publicEntity      *openpgp.Entity
	publicKeyToUse    []*openpgp.Entity
	privateEntityList openpgp.EntityList
	publicKeys        []string
	privateKeys       []string
}

// NoSuchKeyError is thrown when supplied GPG key name not available in private and public keychains
type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return "gogpg: no such key \"" + err.key + "\""
}

// IncorrectPassphrase is thrown when the supplied passphrase doesn't match for the key
type IncorrectPassphrase struct {
	key string
}

func (err IncorrectPassphrase) Error() string {
	return "gogpg: incorrect passphrase for \"" + err.key + "\""
}

// New returns a new GPGStore that can then needs to be initialized with Init()
func New(secretKeyring, publicKeyring string) (*GPGStore, error) {
	gs := new(GPGStore)
	gs.secretKeyring = secretKeyring
	gs.publicKeyring = publicKeyring
	var err error
	gs.publicKeys, err = gs.ListPublicKeys()
	if err != nil {
		return gs, err
	}
	gs.privateKeys, err = gs.ListPrivateKeys()
	if err != nil {
		return gs, err
	}
	return gs, nil
}

// ListPrivateKeys returns a list of the names of keys available in the private key chain
func (gs *GPGStore) ListPrivateKeys() ([]string, error) {
	keyringFileBuffer, err := os.Open(gs.secretKeyring)
	if err != nil {
		panic(err)
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		panic(err)
	}
	keys := []string{}
	for _, key := range entityList {
		for _, id := range key.Identities {
			keys = append(keys, strings.Split(strings.Split(id.Name, " <")[0], " (")[0])
		}
	}
	return keys, nil
}

// ListPublicKeys returns a list of the names of keys available in the public key chain
func (gs *GPGStore) ListPublicKeys() ([]string, error) {
	keyringFileBuffer, err := os.Open(gs.publicKeyring)
	if err != nil {
		panic(err)
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		panic(err)
	}
	keys := []string{}
	for _, key := range entityList {
		for _, id := range key.Identities {
			keys = append(keys, strings.Split(strings.Split(id.Name, " <")[0], " (")[0])
		}
	}
	return keys, nil
}

// Init uses the supplied identity and passphrase to determine the GPG parameters.
// The identity must be in the secret and public keychain, and returns an error otherwise.
// The passphrase is validated and returns an error if doesn't exist.
func (gs *GPGStore) Init(identity, passphrase string) error {
	gs.identity = identity
	gs.passphrase = passphrase

	// Open the private key file
	var entityList []*openpgp.Entity
	keyringFileBuffer, err := os.Open(gs.secretKeyring)
	if err != nil {
		return err
	}
	defer keyringFileBuffer.Close()
	entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return err
	}
	gs.privateKeyToUse = make([]*openpgp.Entity, 1)
	foundKey := false
	for _, key := range entityList {
		for _, id := range key.Identities {
			if strings.Split(strings.Split(id.Name, " <")[0], " (")[0] == identity {
				gs.privateKeyToUse[0] = key
				gs.privateEntity = key
				foundKey = true
				break
			}
		}
	}
	if !foundKey {
		return NoSuchKeyError{identity}
	}
	gs.privateEntityList = entityList

	// Open the public key file
	keyringFileBuffer, err = os.Open(gs.publicKeyring)
	if err != nil {
		return err
	}
	defer keyringFileBuffer.Close()
	entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return err
	}
	gs.publicKeyToUse = make([]*openpgp.Entity, 1)
	foundKey = false
	for _, key := range entityList {
		for _, id := range key.Identities {
			if strings.Split(strings.Split(id.Name, " <")[0], " (")[0] == identity {
				gs.publicKeyToUse[0] = key
				gs.publicEntity = key
				foundKey = true
				break
			}
		}
	}
	if !foundKey {
		return NoSuchKeyError{identity}
	}

	// Test the password
	enc, err := gs.Encrypt([]byte("testing"))
	if err != nil {
		return err
	}
	_, err = gs.Decrypt(enc)
	if err != nil {
		return IncorrectPassphrase{identity}
	}
	return nil
}

// Decrypt uses the supplied GPG parameters to decrypt armor-encoded GPG data.
func (gs *GPGStore) Decrypt(encrypted []byte) (decrypted []byte, err error) {
	passphraseByte := []byte(gs.passphrase)
	gs.privateEntity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range gs.privateEntity.Subkeys {
		err = subkey.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			return
		}
	}
	result, err := armor.Decode(bytes.NewBuffer(encrypted))
	if err != nil {
		return
	}

	md, err := openpgp.ReadMessage(result.Body, gs.privateEntityList, nil, nil)
	if err != nil {
		return
	}
	decrypted, err = ioutil.ReadAll(md.UnverifiedBody)
	return
}

// Encrypt takes a string and returns a string that is armor encoded using the supplied GPG credentials.
func (gs *GPGStore) Encrypt(unencrypted []byte) (encrypted []byte, err error) {
	buf := new(bytes.Buffer)
	msg, _ := armor.Encode(buf, "PGP MESSAGE", nil)
	w, err := openpgp.Encrypt(msg, gs.publicKeyToUse, nil, nil, nil)
	if err != nil {
		return
	}
	_, err = w.Write(unencrypted)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	err = msg.Close()
	if err != nil {
		return
	}

	encrypted, err = ioutil.ReadAll(buf)
	return
}
