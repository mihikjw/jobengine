package queue

import (
	"github.com/MichaelWittgreffe/jobengine/filesystem"
)

//DBFile represents the encrypted file holding persistent queue data
type DBFile struct {
	filepath    string
	fileHandler filesystem.FileSystem
}

//NewDBFile is a constructor for DBFile
func NewDBFile(filepath, fsType string) *DBFile {
	return &DBFile{
		filepath:    filepath,
		fileHandler: filesystem.NewFileSystem(fsType),
	}
}

//Exists checks whether the db file exists at the loaded location or not
func (db *DBFile) Exists() bool {
	exists, _ := db.fileHandler.FileExists(db.filepath)
	return exists
}
