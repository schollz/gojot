package gojot

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type DocCache struct {
	Docs                  Documents
	LoadedFiles           map[string]bool
	DocumentAndEntryNames map[string]map[string]bool
}

func (gj *gojot) SaveDocCache() (err error) {
	cache := DocCache{
		Docs:                  gj.docs,
		LoadedFiles:           gj.loadedFiles,
		DocumentAndEntryNames: gj.documentAndEntryNames,
	}
	bCache, err := json.Marshal(cache)
	if err != nil {
		return
	}
	enc, err := gj.gpg.Encrypt(bCache)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join(gj.root, "cache.json"), enc, 0755)
	return
}

func (gj *gojot) LoadDocCache() (err error) {
	var cache DocCache
	bCache, err := ioutil.ReadFile(path.Join(gj.root, "cache.json"))
	if err != nil {
		return
	}
	dec, err := gj.gpg.Decrypt(bCache)
	if err != nil {
		return
	}
	err = json.Unmarshal(dec, &cache)
	if err != nil {
		return
	}
	gj.docs = cache.Docs
	gj.loadedFiles = cache.LoadedFiles
	gj.documentAndEntryNames = cache.DocumentAndEntryNames
	gj.log.Infof("Loaded %d files from cache", len(gj.loadedFiles))
	return
}
