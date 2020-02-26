package main

import (
	"fmt"
	"os"

	"github.com/MichaelWittgreffe/jobengine/api"
	"github.com/MichaelWittgreffe/jobengine/configload"
	"github.com/MichaelWittgreffe/jobengine/logger"
	"github.com/MichaelWittgreffe/jobengine/queue"
)

func main() {
	logger := logger.NewLogger("std")

	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) <= 0 {
		logger.Info("CONFIG_PATH Not Defined, Using Default")
		configPath = "/jobengine/config.yml"
	}

	dbPath := os.Getenv("DB_PATH")
	if len(dbPath) <= 0 {
		logger.Info("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	cfg, err := configload.LoadConfig(configload.NewConfigLoader(configPath, "os"))
	if err != nil {
		quit(err, logger)
	}

	logger.Info("Config Loaded")
	dbFile := queue.NewDBFile(dbPath, cfg.CryptoSecret, "os")
	var queueCon *queue.Controller

	if dbFile.Exists() {
		logger.Info("DB File Found, Loading...")

		if queueCon, err = queue.NewControllerFromDB(cfg, dbFile); err != nil {
			quit(err, logger)
		}
	} else {
		if queueCon, err = queue.NewControllerFromConfig(cfg); err != nil {
			quit(err, logger)
		}

		dbFile.LoadController(queueCon)
	}

	logger.Info("Queues Loaded")
	logger.Info(fmt.Sprintf("Complete/Failed Jobs Older Than %d Minutes Will Be Deleted\n", cfg.JobKeepMinutes))
	logger.Info(fmt.Sprintf("Jobs Inprogress For %d Minutes Will Be Marked As Failed\n", cfg.JobTimeoutMinutes))

	comms := make(chan bool, 1)
	comms <- false
	go dbFile.Monitor(comms, logger)
	logger.Info("Write Monitor Routine Started")

	server := api.NewAPIServer(cfg.Version, queueCon, comms, logger)
	logger.Info(fmt.Sprintf("API Starting On Port %d\n", cfg.Port))
	err = server.ListenAndServe(cfg.Port)
	quit(err, logger)
}

//quit exits the program with exit code 1 and prints the error if there was one
func quit(err error, logger logger.Logger) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
