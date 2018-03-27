package util

import (
	"testing"

	"fmt"
)

func TestUtil(t *testing.T) {
	paths, _ := IterDir("D:\\tmp\\seaweedfs" )
	fmt.Println(paths)

}
