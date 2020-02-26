package configload

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
)

//HashProcessor implements a hashing processor interface
type HashProcessor interface {
	ProcessFile(filepath string) (string, error)
	Process(in string) (string, error)
}

//NewHashProcessor creates a new instance of a HashProcessor object
func NewHashProcessor(algorithm string) HashProcessor {
	var hasher hash.Hash

	switch {
	case algorithm == "md5":
		hasher = md5.New()
	case algorithm == "sha1":
		hasher = sha1.New()
	case algorithm == "sha256":
		hasher = sha256.New()
	case algorithm == "sha512":
		hasher = sha512.New()
	default:
		return nil
	}

	return &HashProcess{
		hasher: hasher,
		fs:     filesystem.NewFileSystem("os"),
	}
}

//HashProcess implements HashProcessor
type HashProcess struct {
	hasher hash.Hash
	fs     filesystem.FileSystem
}

//ProcessFile runs the configured hashing algorithm on the given file
func (p *HashProcess) ProcessFile(filepath string) (string, error) {
	if len(filepath) <= 0 {
		return "", fmt.Errorf("Invalid Arg")
	}

	exists, err := p.fs.FileExists(filepath)
	if !exists || err != nil {
		return "", fmt.Errorf("File %s Not Found", filepath)
	}

	fileHandle, err := p.fs.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("Error Opening File %s: %s", filepath, err.Error())
	}

	_, err = io.Copy(p.hasher, fileHandle)
	if err != nil {
		return "", fmt.Errorf("Error Setting Up Hasher: %s", err.Error())
	}

	hashBytes := p.hasher.Sum(nil)
	if len(hashBytes) <= 0 {
		return "", fmt.Errorf("Generated Hash Is Empty")
	}

	return hex.EncodeToString(hashBytes), nil
}

//Process runs the configured hashing algorithm on the given string
func (p *HashProcess) Process(in string) (string, error) {
	if len(in) <= 0 {
		return "", fmt.Errorf("Invalid Arg")
	}

	_, err := p.hasher.Write([]byte(in))
	if err != nil {
		return "", fmt.Errorf("Error Writing Hash: %s", err.Error())
	}

	hashBytes := p.hasher.Sum(nil)
	if len(hashBytes) <= 0 {
		return "", fmt.Errorf("Generated Hash Is Empty")
	}

	return hex.EncodeToString(hashBytes), nil
}
