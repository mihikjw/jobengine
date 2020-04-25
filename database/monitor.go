package database

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/logger"
)

// DBMonitor presents an object for write requests and start monitoring the DB
type DBMonitor interface {
	Write()
	Start()
}

// DBFileMonitor is an object responsible for writing changes to the DBFile object to the database
type DBFileMonitor struct {
	dbFile      *DBFile
	dbFilePath  string
	fileHandler DBFileHandler
	flag        chan bool
	log         logger.Logger
}

// NewDBFileMonitor is a constructor for the DBFileMonitor type
func NewDBFileMonitor(dbFile *DBFile, filePath string, dbFileHandler DBFileHandler, logger logger.Logger) DBMonitor {
	if dbFileHandler == nil {
		return nil
	}

	writeFlag := make(chan bool, 1)
	writeFlag <- false

	return &DBFileMonitor{
		dbFile:      dbFile,
		dbFilePath:  filePath,
		fileHandler: dbFileHandler,
		flag:        writeFlag,
		log:         logger,
	}
}

// Write requests the DBFile to be written to file
func (m *DBFileMonitor) Write() {
	m.flag <- true
}

// Start begins the monitoring and write process, blocks current goroutine with infinite loop
func (m *DBFileMonitor) Start() {
	for true {
		select {
		case writeFlag, open := <-m.flag:
			if open {
				if writeFlag {
					if err := m.fileHandler.SaveToFile(m.dbFile, m.dbFilePath); err != nil {
						m.log.Error(fmt.Sprintf("Error Saving DB File: %s", err.Error()))
					}
				}
			} else {
				m.log.Error("DBFile Monitor Write Channel Closed")
			}
		}
	}
}
