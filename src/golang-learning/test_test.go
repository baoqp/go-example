package golang_learning

import (
	"testing"
	"unsafe"
	"fmt"
)

var ready = make(chan int)

type A struct {
	a int
	b string
}

type B struct {
	a int
	b string
	c float32
}


func cast1(a *A) {
	p := unsafe.Pointer(a)
	b := (*B)(p)
	fmt.Println(b)
}

func cast2(b *B) {
	p := unsafe.Pointer(b)
	a := (*A)(p)
	fmt.Println(a)
}




func Test2(t *testing.T) {


	type test struct {
		a int
		s string
	}

	v1 := test{
		a: 5,
		s: "sdaf",
	}

	fmt.Printf("v1: %#v\n", v1)

	b := *(*[unsafe.Sizeof(v1)]byte)(unsafe.Pointer(&v1))
	fmt.Printf("bytes: %#v\n", b)

	v2 := *(*test)(unsafe.Pointer(&b))

	fmt.Printf("v2: %#v\n", v2)
}
