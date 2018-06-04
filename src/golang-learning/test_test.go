package golang_learning

import (
	"testing"
	"fmt"
)
type A struct {
	name string
	next *A
}


func Test6(t *testing.T) {
	a := A{
		name: "A",
	}

	b := A{
		name: "B",
	}

	a.next = &b

	var p **A
	p = &a.next
	*p = &a
	p = &a.next

	fmt.Println(a)
	fmt.Println(p)
	fmt.Println((*p).name)
	fmt.Println((**p).name)


}
