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

func B(a *A) {
	 a.a = "asdasd"
}

func D() []int {
	return []int{1,2,3}
}

func C(arr []int) {
	arr = D()
}

func Test2(t *testing.T) {
	var a  = &A{}
	B(a)
	fmt.Println(a)

	var arr []int
	C(arr)
	fmt.Println(arr)


}