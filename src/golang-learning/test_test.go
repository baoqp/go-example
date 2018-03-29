package golang_learning

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	GetMembers(&sr{})
}

type A struct {
	a string
}

func C(arr *[]int) {
	*arr = []int{1, 1, 1,}
}

func Test2(t *testing.T) {
	var arr = []int{0, 0, 0,}
	C(&arr)
 	fmt.Println(arr)
}
