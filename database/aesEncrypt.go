package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// AESEncrypt handles encryption/decryption for the AES cipher
type AESEncrypt struct {
	key []byte
}

// Encrypt encrypts the given byte slice with the AES cipher
func (e *AESEncrypt) Encrypt(in []byte) ([]byte, error) {
	cypher, err := aes.NewCipher(e.Key())
	if err != nil {
		return nil, fmt.Errorf("Failed To Create AES Key: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(cypher)
	if err != nil {
		return nil, fmt.Errorf("Failed To Create GCM: %s", err.Error())
	}

	cryptoData := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, cryptoData); err != nil {
		return nil, fmt.Errorf("Failed To Gen Random Encryption Data: %s", err.Error())
	}

	encryptedData := gcm.Seal(cryptoData, cryptoData, in, nil)
	return encryptedData, nil
}

// Decrypt decrypts the given byte slice with the AES cipher
func (e *AESEncrypt) Decrypt(in []byte) ([]byte, error) {
	cypher, err := aes.NewCipher(e.Key())
	if err != nil {
		return nil, fmt.Errorf("Failed To Create AES Key: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(cypher)
	if err != nil {
		return nil, fmt.Errorf("Failed To Create GCM: %s", err.Error())
	}

	nonceSize := gcm.NonceSize()
	if len(in) < nonceSize {
		return nil, fmt.Errorf("Failed To Create Nonce Of Required Size, DataSize: %d, NonceSize: %d", len(in), nonceSize)
	}

	nonce, in := in[:nonceSize], in[nonceSize:]
	rawData, err := gcm.Open(nil, nonce, in, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed To Decrypt Data: %s", err.Error())
	}

	return rawData, nil
}

// Key returns the currently set key as []bytes
func (e *AESEncrypt) Key() []byte {
	return e.key
}
