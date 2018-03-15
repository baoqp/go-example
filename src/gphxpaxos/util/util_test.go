package util

import (
	"testing"
	"time"
	"fmt"
)

func TestUtil(t *testing.T) {
	StartRoutine(func(){
		time.Sleep(100 * time.Millisecond)
		fmt.Println("sleep 100 ms done ")
	})
}
