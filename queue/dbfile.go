package queue

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
)

//DBFile represents the encrypted file holding persistent queue data
type DBFile struct {
	filepath     string
	cryptoSecret string
	fileHandler  filesystem.FileSystem
	controller   *Controller
}

//NewDBFile is a constructor for DBFile
func NewDBFile(filepath, cryptoSecret, fsType string) *DBFile {
	return &DBFile{
		filepath:     filepath,
		cryptoSecret: cryptoSecret,
		fileHandler:  filesystem.NewFileSystem(fsType),
		controller:   nil,
	}
}

//Exists checks whether the db file exists at the loaded location or not
func (db *DBFile) Exists() bool {
	exists, _ := db.fileHandler.FileExists(db.filepath)
	return exists
}

//LoadController adds a pointer to the controller to the DBFile
func (db *DBFile) LoadController(in *Controller) {
	db.controller = in
}

//LoadFromFile loads a saved file into the Controller struct
func (db *DBFile) LoadFromFile() error {
	if db.controller == nil {
		return fmt.Errorf("Controller Is Nil")
	}

	exists, err := db.fileHandler.FileExists(db.filepath)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("DB File Not Found At Path %s", db.filepath)
	}

	_, err = db.fileHandler.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("Failed Loading DB File: %s", err.Error())
	}

	//un-encrypt and load data here

	return nil
}

//SaveToFile saves the loaded data onto file
func (db *DBFile) SaveToFile() error {
	if db.controller == nil {
		return fmt.Errorf("Controller Is Nil")
	}

	if exists, err := db.fileHandler.FileExists(db.filepath); exists {
		if err := db.fileHandler.DeleteFile(db.filepath); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	rawData, err := db.controller.ExportQueues()
	if err != nil {
		return fmt.Errorf("Failed To Export Queues: %s", err.Error())
	}

	fileData, err := json.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("Failed To Marshal Queues: %s", err.Error())
	}

	encryptionKey := []byte(db.cryptoSecret)

	cypher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return fmt.Errorf("Failed To Create AES Key: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(cypher)
	if err != nil {
		return fmt.Errorf("Failed To Create GCM: %s", err.Error())
	}

	cryptoData := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, cryptoData); err != nil {
		return fmt.Errorf("Failed To Gen Random Encryption Data: %s", err.Error())
	}

	encryptedData := gcm.Seal(cryptoData, cryptoData, fileData, nil)

	return db.fileHandler.WriteFile(db.filepath, encryptedData)
}

//Monitor is a goroutine to write the DBFile to disk when requested to
func (db *DBFile) Monitor(write chan bool) {
	for true {
		select {
		case writeFlag, open := <-write:
			if open {
				write <- false //reset flag before write

				if writeFlag {
					if err := db.SaveToFile(); err != nil {
						fmt.Printf("Error Saving DB File: %s\n", err.Error())
					}
				}
			} else {
				fmt.Println("DBFile Write Channel Closed")
			}
		}
	}
}
