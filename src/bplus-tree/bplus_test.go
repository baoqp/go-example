package bplus_tree

import (
	"testing"

	"fmt"
)

func Test1(t *testing.T) {
	tree, _ := Open("tree.bp")
	fmt.Println(tree) //&{0 1 64 0 0xc04207c120}
}

func Test2(t *testing.T) {
	tree, _ := Open("tree.bp")
	fmt.Println(tree) //&{0 1 64 0 0xc04207c120}

	for i:=8; i<=10; i++ {
		key := fmt.Sprintf("key %d", i)
		value := fmt.Sprintf("value %d", i)
		tree.Set([]byte(key), []byte(value))
	}

	key := fmt.Sprintf("key %d", 9)
	data, _ := tree.Get([]byte(key))  // TODO 有问题
	fmt.Println(string(data))
}
