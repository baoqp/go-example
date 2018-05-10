package bplus_tree_aof

import (
	"testing"
	"fmt"

	"time"
)

func Test1(t *testing.T) {
	tree, _ := Open("tree.bp")
	fmt.Println(tree) //&{0 1 64 0 0xc04207c120}
}


func TestApi(t *testing.T) {
	tree, _ := Open("tree.bp")

	//start := uint64(time.Now().UnixNano() / 1000000)

	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i)
		tree.Set([]byte(key), []byte(value))
	}

	// first split
	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}
	fmt.Println()


	for i := 7; i <= 8; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i)
		tree.Set([]byte(key), []byte(value))
	}

	for i := 1; i <= 8; i++ {
		key := fmt.Sprintf("%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}
	fmt.Println()

	// second split
	for i := 9; i <= 10; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i)
		tree.Set([]byte(key), []byte(value))
	}

	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}
	fmt.Println()

	key := fmt.Sprintf("%d", 5)
	tree.Remove([]byte(key),nil, nil)
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}

	var rangeValues []string
	var rangeCb RangeCallback = func (arg []byte, key *Key, value *Value) {
		rangeValues = append(rangeValues, string(value.value))
	}
	start := fmt.Sprintf("%d", 1)
	end := fmt.Sprintf("%d", 3)
	tree.GetRange([]byte(start), []byte(end), rangeCb, nil)
	fmt.Println(rangeValues)

}


func TestBench(t *testing.T) {
	tree, _ := Open("tree.bp")

	start := uint64(time.Now().UnixNano() / 1000000)

	n := 1000000

	for i := 1; i <= n; i++ {
		key := fmt.Sprintf("K-%d", i)
		value := fmt.Sprintf("Value-%d", i)
		tree.Set([]byte(key), []byte(value))
	}

	end := uint64(time.Now().UnixNano() / 1000000)

	fmt.Printf("insert %d KV time cost %d ", n, end-start ) // abount 13 seconds
}


func TestCompact(t *testing.T) {

	tree, _ := Open("tree.bp")

	for i := 1000; i <= 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value%d", i)
		tree.Set([]byte(key), []byte(value))
	}


	for i := 1000; i <= 1050; i++ {
		key := fmt.Sprintf("key-%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}
	fmt.Println()

	err := tree.compact()
	if err != nil {
		panic(err)
	}


	for i := 1000; i <= 1050; i++ {
		key := fmt.Sprintf("key-%d", i)
		data, _ := tree.Get([]byte(key))
		fmt.Printf("%s -> %s \n", key, string(data))
	}
	fmt.Println()


}

