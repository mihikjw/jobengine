package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/configload"
)

func main() {
	if _, err := configload.LoadConfig("./examples/config.yml", "os"); err == nil {
		fmt.Println("Config Loaded")
	} else {
		fmt.Println(err.Error())
	}
}
