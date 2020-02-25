package queue

//EncryptionHandler is an interface for encryption/decryption objects
type EncryptionHandler interface {
	Encrypt(in []byte) ([]byte, error)
	Decrypt(in []byte) ([]byte, error)
	Key() []byte
}

//NewEncryptionHandler creates a struct for handler encryption/decryption
func NewEncryptionHandler(key, encType string) EncryptionHandler {
	switch {
	case encType == "AES":
		return &AESEncrypt{
			key: []byte(key),
		}
	default:
		return nil
	}
}
