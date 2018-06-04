package glemon

import (
	"testing"
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


