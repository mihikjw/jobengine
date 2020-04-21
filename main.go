package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/database"
	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/logger"
)

func main() {
	logger := logger.NewLogger("std")
	fileHandler := filesystem.NewFileSystem("os")
	dbFile := database.NewDBFile()

	dbPath := fileHandler.GetEnv("DB_PATH")
	if len(dbPath) <= 0 {
		logger.Info("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	secretKey := fileHandler.GetEnv("SECRET")
	if len(secretKey) <= 0 {
		logger.Fatal("SECRET Not Defined")
	}

	dbFileHandler := database.NewDBFileHandler(
		"fs",
		database.NewEncryptionHandler(secretKey, "AES", database.NewHashHandler("md5")),
		database.NewDBDataHandler("json"),
		fileHandler,
	)

	if dbFileHandler == nil {
		logger.Fatal("Unable To Create File Handler")
	}

	if exists, err := fileHandler.FileExists(dbPath); err == nil {
		if exists {
			if err = dbFileHandler.LoadFromFile(dbFile, dbPath); err == nil {
				logger.Info("Database Loaded")
			}
		} else {
			if err = dbFileHandler.SaveToFile(dbFile, dbPath); err == nil {
				logger.Info("Database Created")
			}
		}
	} else {
		logger.Fatal(fmt.Sprintf("Error Locating DB File: %s", err.Error()))
	}

	dbFileMonitor := database.NewDBFileMonitor(dbFile, dbPath, dbFileHandler, logger)
	if dbFileMonitor == nil {
		logger.Fatal("Failed Creating Monitor")
	}
	go dbFileMonitor.Start()

	// host the API
}
