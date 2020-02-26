package queue

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/logger"
)

//DBFile represents the encrypted file holding persistent queue data
type DBFile struct {
	filepath    string
	crypto      EncryptionHandler
	fileHandler filesystem.FileSystem
	jsonHandler JSONHandler
	controller  *Controller
}

//NewDBFile is a constructor for DBFile
func NewDBFile(filepath, cryptoSecret, fsType string) *DBFile {
	return &DBFile{
		filepath:    filepath,
		crypto:      NewEncryptionHandler(cryptoSecret, "AES"),
		fileHandler: filesystem.NewFileSystem(fsType),
		jsonHandler: new(JSONHandle),
		controller:  nil,
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

	encryptedData, err := db.fileHandler.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("Failed Loading DB File: %s", err.Error())
	}

	rawDataString, err := db.crypto.Decrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("Failed Decrypting DB File: %s", err.Error())
	}

	rawData, err := db.jsonHandler.Unmarshal(rawDataString)
	if err != nil {
		return fmt.Errorf("Failed To Unmarshal Data: %s", err.Error())
	}

	return db.controller.LoadQueues(rawData)
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

	fileData, err := db.jsonHandler.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("Failed To Marshal Queues: %s", err.Error())
	}

	encryptedData, err := db.crypto.Encrypt(fileData)
	if err != nil {
		return fmt.Errorf("Failed to Encrypt Data: %s", err.Error())
	}

	return db.fileHandler.WriteFile(db.filepath, encryptedData)
}

//Monitor is a goroutine to write the DBFile to disk when requested to
func (db *DBFile) Monitor(write chan bool, logger logger.Logger) {
	for true {
		select {
		case writeFlag, open := <-write:
			if open {
				if writeFlag {
					write <- false //reset flag before write

					if err := db.SaveToFile(); err != nil {
						logger.Error(fmt.Sprintf("Error Saving DB File: %s\n", err.Error()))
					}
				}
			} else {
				logger.Error("DBFile Write Channel Closed")
			}
		}
	}
}
