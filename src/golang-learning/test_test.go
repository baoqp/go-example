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

func C(b *B) {
	var len = len(b.as)
	b.as = b.as[:len-1]
}

func Test2(t *testing.T) {
	b := &B{}
	b.as = make([]*A, 0)
	b.as = append(b.as, &A{"a1"})
	fmt.Println(len(b.as))
	C(b)
	fmt.Println(len(b.as))
}
