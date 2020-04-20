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
		logger.Fatal("SECRET Not Defined")
	}

	var err error
	dbFile := database.NewDBFile()
	dbFileHandler := database.NewDBFileHandler(
		"fs",
		database.NewEncryptionHandler(secretKey, "AES", database.NewHashHandler("md5")),
		database.NewDBDataHandler("json"),
	)

	if dbFileHandler == nil {
		logger.Fatal("Unable To Create File Handler")
	}

	if _, err = os.Stat(dbPath); err == nil {
		err = dbFileHandler.LoadFromFile(dbFile, dbPath)
	} else if os.IsNotExist(err) {
		err = dbFileHandler.SaveToFile(dbFile, dbPath)
	} else {
		logger.Error(fmt.Sprintf("Error Locating DB File: %s", err.Error()))
	}

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed To Create Or Load DBFile: %s", err.Error()))
	}

	// start monitoring the DBFile ready for changes
	// host the API
}
