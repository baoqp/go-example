package glemon

import (
	"testing"
	"fmt"
)

func Test1(t *testing.T) {
	main()
}

func Test2(t *testing.T) {

	lemon := &lemon{
		filename: "C:\\test.y",
	}
	Parse(lemon)
}



func Test3(t *testing.T) {
	a := 'A'
	b := byte(a)
	fmt.Printf("%c\n", a)
	fmt.Printf("%c\n", b)
}
