package main

import (
	"fmt"
	"os"

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
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Config Loaded")
	dbFile := queue.NewDBFile(dbPath, "os")
	var queueCon *queue.Controller

	if dbFile.Exists() {
		fmt.Println("DB File Found, Loading...")

		if queueCon, err = queue.NewControllerFromDB(cfg, dbFile); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		if queueCon, err = queue.NewControllerFromConfig(cfg); err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("Queues Loaded")
	queueCon.AddNewQueue("blob", nil)

	/* TO DO
	- add a lock to the dbFile and really understand the priority/how this works in practice
	- spawn a goroutine for writing to the QueueDB & Config
		- requires a channel to take write requests
	- create the API, with access to the queueCon and DB write goroutine request
	*/

	os.Exit(0)
}
