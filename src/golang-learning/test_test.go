package golang_learning

import (
	"testing"
	"fmt"
	"time"
	"container/list"
)

var ready = make(chan int)

func Test2(t *testing.T) {
	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		ready <- 1
	}()
	var i int
	select{
	case i= <- ready:
		fmt.Printf("receive %d", i)
	}

}
