package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/config"
)

func main() {
	if loader := config.NewConfigLoader("/etc/jobengine/config.yml", "os"); loader != nil {
		if _, err := loader.LoadFromFile(); err == nil {
			fmt.Println("Config Loaded")
		} else {
			fmt.Println("Failed To Load Config: " + err.Error())
		}
	} else {
		fmt.Println("Unable To Create Config")
	}
}
