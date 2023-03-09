package hashFS

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"sync"

	"github.com/sirupsen/logrus"
)

type HashFS interface {
	Open(hash string) (fs.File, error)
	GetHash(name string) (string, error)
	GetHashNoErr(name string) string
	PreCache()
}

type hashFS struct {
	hashCache    sync.Map
	reverseCache sync.Map
	backingFS    fs.FS
}

func (hFS *hashFS) Open(input string) (fs.File, error) {
	filename, ok := hFS.hashCache.Load(input)
	if !ok { // file not found
		return nil, fs.ErrInvalid

	} else {
		return hFS.backingFS.Open(filename.(string))
	}
}

func (hFS *hashFS) GetHash(input string) (string, error) {

	if hash, ok := hFS.reverseCache.Load(input); ok {
		return hash.(string), nil
	}
	file, err := hFS.backingFS.Open(input)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hashfn := sha256.New()
	if _, err := io.Copy(hashfn, file); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hashfn.Sum(nil))[0:10] + "." + input
	hFS.hashCache.Store(hash, input)
	hFS.reverseCache.Store(input, hash)
	return hash, nil
}

func (hFS *hashFS) GetHashNoErr(input string) string {
	hash, err := hFS.GetHash(input)
	if err != nil {
		logrus.Panicln(err)
	}
	return hash
}

func GenHashFS(backingFS fs.FS) hashFS {
	return hashFS{
		backingFS: backingFS,
	}
}
