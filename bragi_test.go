package bragi

import (
	"fmt"
	"testing"
)

func TestDebug(t *testing.T) {
	AddError(fmt.Errorf("Some error message")).Debug("Debug info")
	// t.Fatal("not implemented")
}
