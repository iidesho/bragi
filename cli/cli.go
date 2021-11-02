package main

import (
	"fmt"

	log "github.com/Cantara/bragi"
)

func main() {
	fmt.Println("vim-go")
	Close := log.SetOutputFolder("./logs")
	if Close == nil {
		log.Fatal("Could not set logfolder")
	}
	defer Close()
	log.AddError(fmt.Errorf("Some error message")).Debug("Debug info")
	log.Debug("Debug info without error")
}
