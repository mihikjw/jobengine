package main

import (
	"fmt"
	"os"

	"github.com/MichaelWittgreffe/jobengine/api"
	"github.com/MichaelWittgreffe/jobengine/configload"
	"github.com/MichaelWittgreffe/jobengine/queue"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) <= 0 {
		fmt.Println("CONFIG_PATH Not Defined, Using Default")
		configPath = "/jobengine/config.yml"
	}

	dbPath := os.Getenv("DB_PATH")
	if len(dbPath) <= 0 {
		fmt.Println("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	cfg, err := configload.LoadConfig(configload.NewConfigLoader(configPath, "os"))
	if err != nil {
		quit(err)
	}

	fmt.Println("Config Loaded")
	dbFile := queue.NewDBFile(dbPath, cfg.CryptoSecret, "os")
	var queueCon *queue.Controller

	if dbFile.Exists() {
		fmt.Println("DB File Found, Loading...")

		if queueCon, err = queue.NewControllerFromDB(cfg, dbFile); err != nil {
			quit(err)
		}
	} else {
		if queueCon, err = queue.NewControllerFromConfig(cfg); err != nil {
			quit(err)
		}

		dbFile.LoadController(queueCon)
	}

	fmt.Println("Queues Loaded")
	fmt.Printf("Complete/Failed Jobs Older Than %d Minutes Will Be Deleted\n", cfg.JobKeepMinutes)
	fmt.Printf("Jobs Inprogress For %d Minutes Will Be Marked As Failed\n", cfg.JobTimeoutMinutes)

	comms := make(chan bool, 1)
	comms <- false
	go dbFile.Monitor(comms)
	fmt.Println("Write Monitor Routine Started")

	server := api.NewAPIServer(cfg.Version, queueCon, comms)
	fmt.Printf("API Starting On Port %d\n", cfg.Port)
	err = server.ListenAndServe(cfg.Port)
	quit(err)
}

//quit exits the program with exit code 1 and prints the error if there was one
func quit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
