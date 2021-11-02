package main

import (
	"fmt"
	"time"

	log "github.com/cantara/bragi"
)

func main() {
	fmt.Println("vim-go")
	Close := log.SetOutputFolder("./logs")
	if Close == nil {
		log.Fatal("Could not set logfolder")
	}
	defer Close()
	log.StartRotate(nil)
	log.AddError(fmt.Errorf("Some error message")).Debug("Debug info")
	log.Debug("Debug info without error")
	for i := 0; i < 3; i++ {
		time.Sleep(5 * time.Second)
		log.AddError(fmt.Errorf("Some error message")).Warning("Warning info")
	}
}
