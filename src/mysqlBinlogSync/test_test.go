package mysqlBinlogSync

import (
	"testing"
	"fmt"
)

func test(data []byte) {
	data[0] = 0
}

func Test(t *testing.T) {
	//Connect("45.77.120.142:3306", "root", "bqp0205", "dev")
	data := []byte{1, 2, 3, 4}
	test(data)
	fmt.Println(data)
}
