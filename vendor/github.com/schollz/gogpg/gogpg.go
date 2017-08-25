package gogpg

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

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
	logger            *logrus.Logger
	log               *logrus.Entry
}

// NoSuchKeyError is thrown when supplied GPG key name not available in private and public keychains
type NoSuchKeyError struct {
	key string
}

func (err NoSuchKeyError) Error() string {
	return "gogpg: no such key \"" + err.key + "\""
}

// NoSuchKeyRing is thrown when supplied keyring containing the secret and public keys cannot be found
type NoSuchKeyRingError struct {
	key string
}

func (err NoSuchKeyRingError) Error() string {
	return "gogpg: no such key ring \"" + err.key + "\""
}

// IncorrectPassphrase is thrown when the supplied passphrase doesn't match for the key
type IncorrectPassphrase struct {
	key string
}

func (err IncorrectPassphrase) Error() string {
	return "gogpg: incorrect passphrase for \"" + err.key + "\""
}

// New returns a new GPGStore.
// You can specify the verbosity (which uses logrus) and you can optionally specify a keyring folder to use to find secring.gpg and pubring.gpg. If the keyring folder is not supplied, or not found, it will search in the common places to find the correct folder and use that.
func New(debug bool, keyring ...string) (*GPGStore, error) {
	gs := new(GPGStore)
	gs.logger = logrus.New()
	gs.log = gs.logger.WithFields(logrus.Fields{
		"source": "gogpg",
	})
	gs.Debug(debug)

	// check for valid key ring folder
	validKeyringFolder := ""
	if len(keyring) > 0 {
		if exists(keyring[0]) {
			validKeyringFolder = keyring[0]
		}
	}
	// still not found
	for {
		if validKeyringFolder != "" {
			gs.log.Infof("Using keyring '%s'", validKeyringFolder)
			break
		}
		homeDir, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		s := strings.Split(homeDir, `\`)
		s2 := s[len(s)-1]
		s3 := strings.Split(s2, "/")
		s4 := s3[len(s3)-1]
		validKeyringFolder = `C:\cygwin64\home\` + s4 + `\.gnupg`
		if exists(validKeyringFolder) {
			gs.log.Infof("Using keyring '%s'", validKeyringFolder)
			break
		}
		validKeyringFolder = path.Join(homeDir, "AppData/Roaming/gnupg")
		if exists(validKeyringFolder) {
			gs.log.Infof("Using keyring '%s'", validKeyringFolder)
			break
		}
		validKeyringFolder = path.Join(homeDir, ".gnupg")
		if exists(validKeyringFolder) {
			gs.log.Infof("Using keyring '%s'", validKeyringFolder)
			break
		}
		return gs, NoSuchKeyRingError{validKeyringFolder}
	}
	gs.secretKeyring = path.Join(validKeyringFolder, "secring.gpg")
	gs.publicKeyring = path.Join(validKeyringFolder, "pubring.gpg")
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

func (gs *GPGStore) Identity() string {
	return gs.identity
}

func (gs *GPGStore) Debug(on bool) {
	if on {
		gs.logger.SetLevel(logrus.DebugLevel)
	} else {
		gs.logger.SetLevel(logrus.WarnLevel)
	}
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
			idString := strings.Split(strings.Split(id.Name, " <")[0], " (")[0]
			gs.log.Infof("Found id: '%s'", idString)
			keys = append(keys, idString)
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
func (gs *GPGStore) Init(identity string, passphrase string) error {

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

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// BulkDecrypt takes a list of filenames and then decrypts them in parallel and then returns a map with filename keys and the decrypted file contents as values. Files which have errors are quietly skipped.
func (gs *GPGStore) BulkDecrypt(filenames []string, progress ...bool) (filecontents map[string]string, err error) {
	gs.log.Debugf("Decrypting %d files", len(filenames))
	filecontents = make(map[string]string)
	jobs := make(chan string, len(filenames))
	results := make(chan map[string]string, len(filenames))
	for w := 1; w <= runtime.NumCPU(); w++ {
		go gs.worker(w, jobs, results)
	}
	for _, job := range filenames {
		jobs <- job
	}
	close(jobs)
	showProgress := false
	if len(progress) > 0 {
		showProgress = progress[0]
	}
	var bar *pb.ProgressBar
	if showProgress {
		bar = pb.StartNew(len(filenames))
	}
	for i := 0; i < len(filenames); i++ {
		if showProgress {
			bar.Increment()
		}
		data := <-results
		for key := range data {
			if data[key] != "" {
				filecontents[key] = data[key]
			}
		}
	}
	if showProgress {
		bar.Finish()
	}
	gs.log.Debugf("Decrypted %d files", len(filecontents))
	return
}
func (gs *GPGStore) worker(id int, jobs <-chan string, results chan<- map[string]string) {
	for f := range jobs {
		result := make(map[string]string)
		if !exists(f) {
			gs.log.Debugf("Could not find %s", f)
			results <- result
			continue
		}
		data, err := ioutil.ReadFile(f)
		if err != nil {
			gs.log.Debugf("Could not read %s", f)
			results <- result
			continue
		}
		dec, err := gs.Decrypt(data)
		if err != nil {
			gs.log.Debugf("Could not decrypt %s", f)
			results <- result
			continue
		}
		result[f] = string(dec)
		results <- result
	}
}
