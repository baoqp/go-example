package glemon

import (
	"unsafe"
	"testing"
	"fmt"
)

type Node struct {
	v    int
	next *Node
}

func List1() *Node {
	n1 := Node{
		v: 1,
	}
	n3 := Node{
		v: 3,
	}
	n5 := Node{
		v: 5,
	}

	n1.next = &n3
	n3.next = &n5
	return &n1
}

func List2() *Node {
	n2 := Node{
		v: 2,
	}
	n4 := Node{
		v: 4,
	}
	n6 := Node{
		v: 6,
	}

	n2.next = &n4
	n4.next = &n6
	return &n2
}

// type comparator func(pointer unsafe.Pointer, pointer2 unsafe.Pointer) int
func cmpNode(pointer unsafe.Pointer, pointer2 unsafe.Pointer) int {
	n1 := (*Node)(pointer)
	n2 := (*Node)(pointer2)
	return n1.v - n2.v
}
func getNextNode(p unsafe.Pointer) unsafe.Pointer {
	n := (*Node)(p)
	return unsafe.Pointer(n.next)
}

func setNextNode(p unsafe.Pointer, m unsafe.Pointer) {
	n := (*Node)(p)
	next := (*Node)(m)
	n.next = next
}

func TestMerge(t *testing.T) {
	n1 := List1()
	n2 := List2()

	headP := merge(unsafe.Pointer(n1), unsafe.Pointer(n2), cmpNode, getNextNode, setNextNode)
	head := (*Node)(headP)
	for head != nil {
		fmt.Println(head.v)
		head = head.next
	}
}

func TestMsort(t *testing.T) {
	n1 := List1()
	n2 := List2()
	head := n1
	ptr := head
	for n1.next != nil {
		n1 = n1.next
	}
	n1.next = n2
	fmt.Println("before sort")
	for ptr != nil {
		fmt.Printf("%d, ", ptr.v)
		ptr = ptr.next
	}

	headP := msort(unsafe.Pointer(head), cmpNode, getNextNode, setNextNode)
	head = (*Node)(headP)
	fmt.Println("after sort")
	for head != nil {
		fmt.Printf("%d, ", head.v)
		head = head.next
	}
}

