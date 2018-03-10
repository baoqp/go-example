package graft

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Println(ParseUri("http://www.baidu.com?a=1"))
}
