package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/pkg/api"
	"github.com/MichaelWittgreffe/jobengine/pkg/crypto"
	"github.com/MichaelWittgreffe/jobengine/pkg/database"
	"github.com/MichaelWittgreffe/jobengine/pkg/filesystem"
	"github.com/MichaelWittgreffe/jobengine/pkg/logger"
)

func main() {
	logger := logger.NewLogger("std")
	fileHandler := filesystem.NewFileSystem("os")
	dbFile := database.NewDBFile()
	dbPath, apiPort, secretKey := getEnvVars(logger, fileHandler)

	dbFileHandler := database.NewDBFileHandler(
		"fs",
		crypto.NewEncryptionHandler(secretKey, "AES", crypto.NewHashHandler("md5")),
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

	httpAPI := api.NewHTTPAPI(logger, dbFileMonitor, database.NewQueryController(dbFile, crypto.NewHashHandler("sha512")))
	logger.Info(fmt.Sprintf("Started Listening On Port %s", apiPort))
	logger.Fatal(httpAPI.ListenAndServe(apiPort).Error())
}

// getEnvVars returns the required env var values/default values - exits app if mandatory values are not populated
func getEnvVars(l logger.Logger, fh filesystem.FileSystem) (string, string, string) {
	dbPath := fh.GetEnv("DB_PATH")
	if len(dbPath) <= 0 {
		l.Info("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	apiPort := fh.GetEnv("API_PORT")
	if len(dbPath) <= 0 {
		l.Info("API_PORT Not Defined, Using Default")
		apiPort = "80"
	}

	secretKey := fh.GetEnv("SECRET")
	if len(secretKey) <= 0 {
		l.Fatal("SECRET Not Defined")
	}

	return dbPath, apiPort, secretKey
}
