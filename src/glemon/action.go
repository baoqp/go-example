package glemon

import "unsafe"

var action_freelist []*action

func Action_new() *action {
	var new *action
	if action_freelist == nil || len(action_freelist) == 0 { // TODO 其实是一个对象池的概念
		amt := 100
		action_freelist = make([]*action, 0)
		for i := 0; i < amt; i++ {
			action_freelist = append(action_freelist, &action{})
		}

		for i := 0; i < amt-1; i++ {
			action_freelist[i].next = action_freelist[i+1]
		}
	}
	new = action_freelist[0]
	action_freelist = action_freelist[1:]
	return new
}

// Compare two actions for sorting purposes.  Return negative, zero, or positive
// if the first action is less than, equal to, or greater than the first
func actioncmp(ap1 *action, ap2 *action) int {
	var rc int
	rc = ap1.sp.index - ap2.sp.index
	if rc == 0 {
		rc = int(ap1.typ - ap2.typ)
	}

	if rc == 0 {
		//assert( ap1.type==REDUCE || ap1.type==RD_RESOLVED || ap1.type==CONFLICT);
		//assert( ap2.type==REDUCE || ap2.type==RD_RESOLVED || ap2.type==CONFLICT);
		rc = ap1.rp.index - ap2.rp.index
	}

	return rc
}

func cmpAction(pointer unsafe.Pointer, pointer2 unsafe.Pointer) int {
	a1 := (*action)(pointer)
	a2 := (*action)(pointer2)
	return actioncmp(a1, a2)
}
func getNextAction(p unsafe.Pointer) unsafe.Pointer {
	a := (*action)(p)
	return unsafe.Pointer(a.next)
}

func setNextAction(p unsafe.Pointer, m unsafe.Pointer) {
	n := (*action)(p)
	next := (*action)(m)
	n.next = next
}

func Action_sort(ap *action) *action {
	return (*action)(msort(unsafe.Pointer(ap), cmpAction, getNextAction, setNextAction))
}

func Action_add(app **action, typ e_action, sp *symbol, arg unsafe.Pointer) {
	new := Action_new()
	new.next = *app
	*app = new
	new.typ = typ
	new.sp = sp
	if typ == SHIFT {
		new.stp = (*state)(arg)
	} else {
		new.rp = (*rule)(arg)
	}

}
