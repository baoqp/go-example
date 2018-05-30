package golang_learning

import (
	"testing"
	"unsafe"
	"fmt"
)

var ready = make(chan int)

type A struct {
	a int
	b string
}

type B struct {
	a int
	b string
	c float32
}

func cast1(a *A) {
	p := unsafe.Pointer(a)
	b := (*B)(p)
	fmt.Println(b)
}

func cast2(b *B) {
	p := unsafe.Pointer(b)
	a := (*A)(p)
	fmt.Println(a)
}

func Test2(t *testing.T) {

	type test struct {
		a int
		s string
	}

	v1 := test{
		a: 5,
		s: "sdaf",
	}

	fmt.Printf("v1: %#v\n", v1)

	b := *(*[unsafe.Sizeof(v1)]byte)(unsafe.Pointer(&v1))
	fmt.Printf("bytes: %#v\n", b)

	v2 := *(*test)(unsafe.Pointer(&b))

	fmt.Printf("v2: %#v\n", v2)
}

type pageId uint32
type DataPage struct {
	pageNo pageId
	total  uint16 // 写入次数
	curr   uint16 // data当前索引
	dirty  bool
	data   [4000]byte
}

func Test3(t *testing.T) {
	fmt.Println(unsafe.Sizeof(DataPage{}))
}




func strhash(x string) int {
	h := 0
	for _, c := range x {
		h = h*13 + int(c)
	}
	return h
}

type s_x3node struct {
	data string
	key  string
}


type s_x3 struct {
	size  int
	count int
	tbl   []*s_x3node
	ht    [][]*s_x3node
}

var x3a *s_x3

func  State_init() {
	if x3a != nil {
		return
	}

	x3a = &s_x3{}

	x3a.size = 8
	x3a.count = 0
	x3a.tbl = make([]*s_x3node, 0, 8)
	x3a.ht = make([][]*s_x3node, 8, 8)

}

// To be test
func  State_insert(data string, key string) bool {

	if x3a == nil {
		return false
	}

	ph := strhash(key)
	h := ph & (x3a.size - 1)
	np := x3a.ht[h]

	for _, node := range np {
		if node.key == key {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return false
		}
	}

	// 扩容
	if x3a.count == x3a.size {
		size := x3a.size * 2
		arr := make([]*s_x3node, x3a.size, size)
		copy(arr, x3a.tbl)
		x3a.ht = make([][]*s_x3node, size, size)
		for i := 0; i < x3a.count; i++ {
			h = strhash(x3a.tbl[i].key) & (size - 1)
			x3a.ht[h] = append(x3a.ht[h], x3a.tbl[i])
		}

		x3a.tbl = arr
		x3a.size = size
	}

	newstate := &s_x3node{
		key:  key,
		data: data,
	}
	x3a.tbl = append(x3a.tbl, newstate)
	h = ph & (x3a.size - 1)
	x3a.ht[h] = append(x3a.ht[h], newstate)
	x3a.count += 1
	return true
}

func  State_find(key string) string {
	if x3a == nil {
		return ""
	}


	h := strhash(key)  & (x3a.size - 1)

	np := x3a.ht[h]

	for _, node := range np {

		if node.key ==  key  {
			return node.data
		}
	}
	return ""
}


func Test4(t *testing.T) {
	State_init()

	for i:=0; i<20; i++ {

		key := fmt.Sprintf("key-%d", i)
		data := fmt.Sprintf("data-%d", i)
		State_insert(data, key)
		data_ := State_find(key)
		fmt.Printf("%s <-> %s \n", data, data_)
	}
}

func Test5(t *testing.T) {
	a := []int{1,2,3}
	b := make([]int, len(a), 10)
	copy(b, a)
	fmt.Println(b)
}