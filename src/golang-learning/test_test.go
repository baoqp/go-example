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



func Test2(t *testing.T) {
	a := &A{"a"}
	C(a)
	fmt.Println(a)
}
