package database

import (
	"fmt"
	"io/ioutil"
)

// DBFileHandler provides an interface for an object to load/save the db file to disk
type DBFileHandler interface {
	SaveToFile(dbFile *DBFile, filePath string) error
	LoadFromFile(dbFile *DBFile, filePath string) error
}

// NewDBFileHandler is a factory function for creating a derived instance of the DBFileHandler interface. Performs a hash on the given key to ensure size
func NewDBFileHandler(dbFileHandleType string, encryptHandler EncryptionHandler, dataHandler DBDataHandler) DBFileHandler {
	if len(dbFileHandleType) == 0 || encryptHandler == nil || dataHandler == nil {
		return nil
	}

	switch {
	case dbFileHandleType == "fs":
		return &FSFileHandler{
			encrypt: encryptHandler,
			data:    dataHandler,
		}
	default:
		return nil
	}
}

// FSFileHandler handles a DBFile through the default avalible file system
type FSFileHandler struct {
	encrypt EncryptionHandler
	data    DBDataHandler
}

// SaveToFile saves the given dbFile to the given filePath
func (h *FSFileHandler) SaveToFile(dbFile *DBFile, filePath string) error {
	if dbFile == nil || len(filePath) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	dbFile.mutex.Lock()
	defer dbFile.mutex.Unlock()

	encodedData, err := h.data.Encode(dbFile)
	if err != nil {
		return err
	}

	encryptedData, err := h.encrypt.Encrypt(encodedData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("Failed To Save DBFile: %s", err.Error())
	}

	return nil
}

// LoadFromFile loads the given filePath database into the given dbFile object
func (h *FSFileHandler) LoadFromFile(dbFile *DBFile, filePath string) error {
	if dbFile == nil || len(filePath) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	dbFile.mutex.Lock()
	defer dbFile.mutex.Unlock()

	/* ADD THE IMPLEMENTATION OF DATABASE LOADING HERE */

	return fmt.Errorf("Not Implemented")
}
