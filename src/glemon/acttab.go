package glemon

type yyaction struct {
	lookahead int // Value of the lookahead token
	action    int // Action to take on the given lookahead
}

type acttab struct {
	nAction         int         // Number of used slots in aAction[]
	nActionAlloc    int         // Slots allocated for aAction[]
	aAction         []*yyaction // The yy_action[] table under construction
	aLookahead      []*yyaction // A single new transaction set
	mnLookahead     int         // Minimum aLookahead[].lookahead
	mnAction        int         // Action associated with mnLookahead
	mxLookahead     int         // Maximum aLookahead[].lookahead
	nLookahead      int         // Used slots in aLookahead[]
	nLookaheadAlloc int         // Slots allocated in aLookahead[]
}

func acttab_size(a *acttab) int {
	return a.nAction
}

func acttab_yyaction(a *acttab, N int) int {
	return a.aAction[N].action
}

func acttab_yylookahead(a *acttab, N int) int {
	return a.aAction[N].lookahead
}

func acttab_free(a *acttab) {
	//donothing
}

func acttab_alloc() *acttab {
	return &acttab{}
}

// Add a new action to the current transaction set
func acttab_action(p *acttab, lookahead int, action int) {
	if p.nLookahead == 0 {
		p.mxLookahead = lookahead
		p.mnLookahead = lookahead
		p.mnAction = action
	} else {
		if p.mxLookahead < lookahead {
			p.mxLookahead = lookahead
		}
		if p.mnLookahead > lookahead {
			p.mnLookahead = lookahead
			p.mnAction = action
		}
	}
	p.aLookahead = append(p.aLookahead, &yyaction{lookahead: lookahead, action: action})
	p.nLookahead++
}

//
// Add the transaction set built up with prior calls to acttab_action()
// into the current action table.  Then reset the transaction set back
// to an empty set in preparation for a new round of acttab_action() calls.
//
// Return the offset into the action table of the new transaction.
// TODO
func acttab_insert(p *acttab) int {
	var i, j, k, n int
	// assert( p.nLookahead > 0 )

	// Scan the existing action table looking for an offset where we can
	// insert the current transaction set.  Fall out of the loop when that
	// offset is found.  In the worst case, we fall out of the loop when
	// i reaches p.nAction, which means we append the new transaction set.
	//
	// i is the index in p.aAction[] where p.mnLookahead is inserted.
	for i = 0; i < p.nAction+p.mnLookahead; i++ {
		if p.aAction[i].lookahead < 0 {
			for j = 0; j < p.nLookahead; j++ {
				k = p.aLookahead[j].lookahead - p.mnLookahead + i
				if k < 0 {
					break
				}
				if p.aAction[k].lookahead >= 0 {
					break
				}
			}
			if ( j < p.nLookahead ) {
				continue
			}

			for j = 0; j < p.nAction; j++ {
				if p.aAction[j].lookahead == j+p.mnLookahead-i {
					break
				}
			}
			if j == p.nAction {
				break; /* Fits in empty slots */
			}
		} else if p.aAction[i].lookahead == p.mnLookahead {
			if p.aAction[i].action != p.mnAction {
				continue
			}
			for j = 0; j < p.nLookahead; j++ {
				k = p.aLookahead[j].lookahead - p.mnLookahead + i
				if k < 0 || k >= p.nAction {
					break
				}
				if p.aLookahead[j].lookahead != p.aAction[k].lookahead {
					break
				}
				if p.aLookahead[j].action != p.aAction[k].action {
					break
				}
			}
			if j < p.nLookahead {
				continue
			}
			n = 0
			for j = 0; j < p.nAction; j++ {
				if p.aAction[j].lookahead < 0 {
					continue
				}
				if p.aAction[j].lookahead == j+p.mnLookahead-i {
					n++
				}
			}
			if n == p.nLookahead {
				break; /* Same as a prior transaction set */
			}
		}
	}

	// Insert transaction set at index i.
	for j = 0; j < p.nLookahead; j++ {
		k = p.aLookahead[j].lookahead - p.mnLookahead + i
		p.aAction[k] = p.aLookahead[j]
		if k >= p.nAction {
			p.nAction = k + 1
		}
	}
	p.nLookahead = 0

	// Return the offset that is added to the lookahead in order to get the
	// index into yy_action of the action
	return i - p.mnLookahead;
}
