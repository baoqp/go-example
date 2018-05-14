package gushroom

import (
	"testing"
	"fmt"
)

func A(data []byte) {
	data[0] = 32
	fmt.Println(data)
}

func Test(t *testing.T) {
	data := []byte{1, 2, 3}
	A(data)
	fmt.Println(data)

}

func Test2(t *testing.T) {

}
