package main

import (
	"fmt"
	"os"

	"github.com/MichaelWittgreffe/jobengine/database"
	"github.com/MichaelWittgreffe/jobengine/logger"
)

func main() {
	logger := logger.NewLogger("std")

	dbPath := os.Getenv("DB_PATH")
	if len(dbPath) <= 0 {
		logger.Info("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	secretKey := os.Getenv("SECRET")
	if len(secretKey) <= 0 {
		logger.Error("SECRET Not Defined")
	}

	var err error
	dbFile := database.NewDBFile()
	dbFileHandler := database.NewDBFileHandler(
		"fs",
		database.NewEncryptionHandler(secretKey, "AES", database.NewHashHandler("md5")),
		database.NewDBDataHandler("json"),
	)

	if dbFileHandler == nil {
		logger.Error("Unable To Create File Handler")
		os.Exit(1)
	}

	if _, err = os.Stat(dbPath); err == nil {
		err = dbFileHandler.LoadFromFile(dbFile, dbPath)
	} else if os.IsNotExist(err) {
		err = dbFileHandler.SaveToFile(dbFile, dbPath)
	} else {
		logger.Error(fmt.Sprintf("Error Locating DB File: %s", err.Error()))
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Failed To Create Or Load DBFile: %s", err.Error()))
		os.Exit(1)
	}

	// start monitoring the DBFile ready for changes
	// host the API
}
