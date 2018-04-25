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


	keys := tree.header.page.keys

	fmt.Println(keys)

	//start := uint64(time.Now().UnixNano() / 1000000)

	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i)
		tree.Set([]byte(key), []byte(value))
	}

	// TODO 分裂之后，在插入并没有保存到keys里面或者没有保存到文件中。。没看出来错在哪里
	for i := 7; i <= 8; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i)
		tree.Set([]byte(key), []byte(value))
	}


	for i := 9; i <= 10; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("value#%d", i)
		keybyte := []byte(key)
		valuebyte := []byte(value)
		tree.Set(keybyte, valuebyte)
	}

	key := fmt.Sprintf("%d", 3)
	data, _ := tree.Get([]byte(key))

	fmt.Printf("%s -> %s", key, string(data))
}

// 最后一个header write head &{3264 743 64 859464750716865948 0xc04207e120}
