package database

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

// HashHandler implements a hashing processor interface
type HashHandler interface {
	Process(input string) (string, error)
}

// NewHashHandler creates a new instance of a HashHandler object
func NewHashHandler(algorithm string) HashHandler {
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

	return &HashProcess{hasher: hasher}
}

// HashProcess is an object for performing hash functions on strings
type HashProcess struct {
	hasher hash.Hash
}

// Process performs the configured hash function on the given input string
func (hs *HashProcess) Process(input string) (string, error) {
	if len(input) == 0 {
		return "", fmt.Errorf("Invalid Arg")
	}

	_, err := hs.hasher.Write([]byte(input))
	if err != nil {
		return "", fmt.Errorf("Error Writing Hash: %s", err.Error())
	}

	hashBytes := hs.hasher.Sum(nil)
	if len(hashBytes) <= 0 {
		return "", fmt.Errorf("Generated Hash Is Empty")
	}

	return hex.EncodeToString(hashBytes), nil
}
