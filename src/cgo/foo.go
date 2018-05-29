package main


// #include <stdio.h>
// #include <stdlib.h>
// #include "foo.h"
import "C"
import "fmt"

// 参考 https://www.jianshu.com/p/ce97accb1801
// go build
// cmd下运行cgo
func main() {
	fmt.Println(C.count)
	C.foo()
}
