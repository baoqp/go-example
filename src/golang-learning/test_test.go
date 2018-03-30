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


type B struct {
	as []*A
}

func C(a *A) {
	a1 := &A{"a1"}
	*a = *a1
}

func D() {
	var a = make(map[int]string)
	a[1] = "1"
	a[2] = "2"


	for k := range a {
		if k == 1 {
			delete(a, 1)
		}
	}

	fmt.Println(a)
}

func Test2(t *testing.T) {
	D()
}
