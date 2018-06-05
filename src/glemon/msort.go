package glemon

import "unsafe"

/*
** A generic merge-sort program.
**
** USAGE:
** Let "ptr" be a pointer to some structure which is at the head of
** a null-terminated list.  Then to sort the list call:
**
**     ptr = msort(ptr,&(ptr->next),cmpfnc);
**
** In the above, "cmpfnc" is a pointer to a function which compares
** two instances of the structure and returns an integer, as in
** strcmp.  The second argument is a pointer to the pointer to the
** second element of the linked list.  This address is used to compute
** the offset to the "next" field within the structure.  The offset to
** the "next" field must be constant for all structures in the list.
**
** The function returns a new pointer which is the head of the list
** after sorting.
**
** ALGORITHM:
** Merge-sort.
*/

type comparator func(pointer unsafe.Pointer, pointer2 unsafe.Pointer) int
type getNext func(p unsafe.Pointer) unsafe.Pointer
type setNext func(p unsafe.Pointer, m unsafe.Pointer)

// merge two sorted linked list 合并两个已经排序的链表
func merge(a unsafe.Pointer, b unsafe.Pointer, cmp comparator, getNext getNext, setNext setNext) unsafe.Pointer {
	var head unsafe.Pointer;
	var ptr unsafe.Pointer;
	if a == nil {
		head = b
	} else if b == nil {
		head = a
	} else {
		if cmp(a, b) <= 0 {
			ptr = a
			a = getNext(a)
		} else {
			ptr = b;
			b = getNext(b)
		}
		head = ptr
		for a != nil && b != nil {
			if cmp(a, b) <= 0 {
				setNext(ptr, a)
				ptr = a
				a = getNext(a)
			} else {
				setNext(ptr, b)
				ptr = b
				b = getNext(b)
			}
		}

		if a != nil {
			setNext(ptr, a)
		} else {
			setNext(ptr, b)
		}

	}

	return head
}

const LISTSIZE = 30

// TODO implment a generic msort gracefully
func msort(list unsafe.Pointer, cmp comparator, getNext getNext, setNext setNext) unsafe.Pointer {
	var ep unsafe.Pointer
	var set [LISTSIZE]unsafe.Pointer
	var i int
	for list != nil {
		ep = list
		list = getNext(list)
		setNext(ep, nil)
		// 如果s[i]=1，即挂载了一条链表，那么这条链表的长度为2^i
		for i = 0; i < LISTSIZE-1 && set[i] != nil; i++ {
			ep = merge(ep, set[i], cmp, getNext, setNext)
			set[i] = nil
		}
		set[i] = ep
	}
	ep = nil
	for i=0;i<LISTSIZE; i++ {
		if set[i] != nil {
			ep = merge(set[i], ep, cmp, getNext, setNext)
		}
	}

	return ep
}
