package main

import (
	"fmt"
	"os"

	"github.com/MichaelWittgreffe/jobengine/api"
	"github.com/MichaelWittgreffe/jobengine/configload"
	"github.com/MichaelWittgreffe/jobengine/models"
	"github.com/MichaelWittgreffe/jobengine/queue"
)

const configPath string = "./examples/config.yml"
const dbPath string = "./examples/queues.queuedb"

func main() {
	var err error
	var cfg *models.Config

	if cfg, err = configload.LoadConfig(configPath, "os"); err != nil {
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
	fmt.Printf("Jobs Older Than %d Minutes Will Be Deleted\n", cfg.JobKeepMinutes)

	comms := make(chan bool, 1)
	comms <- false
	go dbFile.Monitor(comms)
	fmt.Println("Write Monitor Routine Started")

	server := api.NewAPIServer(cfg.Version, queueCon, comms)
	fmt.Printf("API Starting On Port %d\n", cfg.Port)
	err = server.ListenAndServe(cfg.Port)
	quit(err)
}

//quit exits the program with an exit code and prints the error if there was one
func quit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
