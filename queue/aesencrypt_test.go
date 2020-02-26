package queue

import (
	"testing"
)

//TestEncryptDecrypt1 succesfully encrypts and decrypts a string
func TestEncryptDecrypt1(t *testing.T) {
	input := "hello_world"
	encryptor := NewEncryptionHandler("32bytetestkey_g69j28hdbi='a3[f]4", "AES")

	encryptedResult, err := encryptor.Encrypt([]byte(input))
	if err != nil {
		t.Errorf("TestEncryptDecrypt1: Unexpected Encrypt Error: %s", err.Error())
	}
	if encryptedResult == nil {
		t.Errorf("TestEncryptDecrypt1: Encrypt Result Unexpectedly Nil")
	}

	decryptedResult, err := encryptor.Decrypt(encryptedResult)
	if err != nil {
		t.Errorf("TestEncryptDecrypt1: Unexpected Decrypt Error: %s", err.Error())
	}
	if decryptedResult == nil {
		t.Errorf("TestEncryptDecrypt1: Decrypt Result Unexpectedly Nil")
	}

	output := string(decryptedResult)

	if output != input {
		t.Errorf("TestEncryptDecrypt1: Input & Output Mismatch: %s != %s", input, output)
	}
}
