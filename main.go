package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/config"
)

func main() {
	if cfg := config.NewConfig("/etc/jobengine/config.json", "os"); cfg != nil {
		if err := cfg.LoadFromFile(); err == nil {
			fmt.Println("Config Loaded")
		} else {
			fmt.Println("Failed To Load Config: " + err.Error())
		}
	} else {
		fmt.Println("Unable To Create Config")
	}
}
