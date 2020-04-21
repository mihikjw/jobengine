package database

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
)

// DBFileHandler provides an interface for an object to load/save the db file to disk
type DBFileHandler interface {
	SaveToFile(dbFile *DBFile, filePath string) error
	LoadFromFile(dbFile *DBFile, filePath string) error
}

// NewDBFileHandler is a factory function for creating a derived instance of the DBFileHandler interface. Performs a hash on the given key to ensure size
func NewDBFileHandler(dbFileHandleType string, encryptHandler EncryptionHandler, dataHandler DBDataHandler, fileHandler filesystem.FileSystem) DBFileHandler {
	if len(dbFileHandleType) == 0 || encryptHandler == nil || dataHandler == nil {
		return nil
	}

	switch {
	case dbFileHandleType == "fs":
		return &FSFileHandler{
			encrypt: encryptHandler,
			data:    dataHandler,
			file:    fileHandler,
		}
	default:
		return nil
	}
}

// FSFileHandler handles a DBFile through the default avalible file system
type FSFileHandler struct {
	encrypt EncryptionHandler
	data    DBDataHandler
	file    filesystem.FileSystem
}

// SaveToFile saves the given dbFile to the given filePath, applies lock
func (h *FSFileHandler) SaveToFile(dbFile *DBFile, filePath string) error {
	if dbFile == nil || len(filePath) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	dbFile.lock.Lock()
	defer dbFile.lock.Unlock()

	if exists, err := h.file.FileExists(filePath); exists && err == nil {
		if err = h.file.DeleteFile(filePath); err != nil {
			return fmt.Errorf("Error Deleting Existing DBFile: %s", err.Error())
		}
	} else if err != nil {
		return fmt.Errorf("Error Checking DBFile Existence: %s", err.Error())
	}

	encodedData, err := h.data.Encode(dbFile)
	if err != nil {
		return fmt.Errorf("Failed Encoding Data: %s", err)
	}

	encryptedData, err := h.encrypt.Encrypt(encodedData)
	if err != nil {
		return fmt.Errorf("Failed Encrypting Data: %s", err)
	}

	if err = h.file.WriteFile(filePath, encryptedData); err != nil {
		return fmt.Errorf("Failed To Save DBFile: %s", err.Error())
	}

	return nil
}

// LoadFromFile loads the given filePath database into the given dbFile object, applies lock
func (h *FSFileHandler) LoadFromFile(dbFile *DBFile, filePath string) error {
	if dbFile == nil || len(filePath) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	dbFile.lock.Lock()
	defer dbFile.lock.Unlock()

	if exists, err := h.file.FileExists(filePath); !exists && err == nil {
		return fmt.Errorf("DBFile Not Found At Path %s", filePath)
	} else if err != nil {
		return fmt.Errorf("Error Checking DBFile Existence: %s", err.Error())
	}

	encryptedData, err := h.file.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error Reading Data: %s", err)
	}

	decryptedData, err := h.encrypt.Decrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("Error Decrypting Data: %s", err.Error())
	}

	err = h.data.Decode(decryptedData, dbFile)
	if err != nil {
		return fmt.Errorf("Error Decoding Data: %s", err.Error())
	}

	return nil
}
