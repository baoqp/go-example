package glemon

/*
** Routines processing configuration follow-set propagation links
** in the LEMON parser generator.
** propagation links TODO ???
*/

var plink_freelist []*plink

func Plink_new() *plink {
	var new *plink
	if plink_freelist == nil || len(plink_freelist) == 0 {
		amt := 100
		plink_freelist = make([]*plink, 0)
		for i := 0; i < amt; i++ {
			plink_freelist = append(plink_freelist, &plink{})
		}
		for i := 0; i < amt-1; i++ {
			plink_freelist[i].next = plink_freelist[i+1]
		}
	}
	new = plink_freelist[0]
	freelist = freelist[1:]
	return new
}

// TODO
func Plink_add(plpp **plink, cfp *config) {
	new := Plink_new()
	new.next = *plpp // 新节点作为链表首节点
	*plpp = new
	new.cfp = cfp
}

// Transfer every plink on the list "from" to the list "to"
// 把from反转接到to上面
func Plink_copy(to **plink, from *plink) {
	var nextpl *plink
	for from != nil {
		nextpl = from.next
		from.next = *to
		*to = from
		from = nextpl
	}
}


// Delete every plink on the list 把以plp开始的plink都归还给plink_freelist
func Plink_delete(plp *plink) {
	var nextpl *plink
	for plp != nil {
		nextpl = plp.next
		plp.next = plink_freelist[0]
		plink_freelist = append([]*plink{plp}, plink_freelist...)
		plp = nextpl;
	}
}

