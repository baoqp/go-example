package mysqlBinlogSync

import (
	"testing"
	"fmt"
)


func Test(t *testing.T) {
	_, err := Connect("45.77.120.142:3306", "root", "bqp0205", "dev")
	if err != nil {
		fmt.Println(err)
	}
}
