package bragi

import (
	"fmt"
	"testing"
	"time"
)

func TestDebug(t *testing.T) {
	SetLevel(INFO)
	AddError(fmt.Errorf("Some error message")).Debug("Debug info")
	// t.Fatal("not implemented")
}

func TestLongRunning(t *testing.T) {
	return
	SetPrefix("bragi_test")
	c := SetOutputFolder("./log")
	defer c()
	i := 0
	ticker := time.Tick(time.Second)
	for range ticker {
		AddError(fmt.Errorf("%d", i)).Info("text")
		i++
	}
}
