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
	p.nLookahead++;
}


//
// Add the transaction set built up with prior calls to acttab_action()
// into the current action table.  Then reset the transaction set back
// to an empty set in preparation for a new round of acttab_action() calls.
//
// Return the offset into the action table of the new transaction.
// TODO
func acttab_insert(p *acttab) {

}
