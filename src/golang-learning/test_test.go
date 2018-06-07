package golang_learning

import (
	"testing"
	"fmt"
	"unicode"
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


func Test7(t *testing.T) {
	cp := "   hello world"
	for len(cp) >0 &&  unicode.IsSpace(rune(cp[0])) {

			cp = cp[1:]

	}
	fmt.Println(cp)
}

