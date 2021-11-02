package main

import (
	log "bragi"
	"fmt"
)

func main() {
	fmt.Println("vim-go")
	log.AddError(fmt.Errorf("Some error message")).Debug("Debug info")
	log.Debug("Debug info without error")
}
