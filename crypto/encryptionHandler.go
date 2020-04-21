package crypto

// EncryptionHandler is an interface for encryption/decryption objects
type EncryptionHandler interface {
	Encrypt(in []byte) ([]byte, error)
	Decrypt(in []byte) ([]byte, error)
	Key() []byte
}

// NewEncryptionHandler creates a struct for handler encryption/decryption
func NewEncryptionHandler(key, encType string, hashHandler HashHandler) EncryptionHandler {
	if len(key) == 0 || len(encType) == 0 || hashHandler == nil {
		return nil
	}

	hashedKey, err := hashHandler.Process(key)
	if err != nil {
		return nil
	}

	switch {
	case encType == "AES":
		return &AESEncrypt{
			key: []byte(hashedKey),
		}
	default:
		return nil
	}
}
