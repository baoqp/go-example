package glemon

import (
	"unsafe"
	"fmt"
)

/*
** Routines to processing a configuration list and building a state
** in the LEMON parser generator.
** config是其他语法书中的LR item
** 一个状态包含的所有项目，分为两种类型：
** 1.基本项目（basis configuration），或者称为核心项目（kernel configuration），是指初始项以及所有分割点不在最左端的项目
** 2.非基本项目，所有分割点在最左端的非初始项目，即可以用该项目的产生式表示，而省略0处的分割点。所有的非基本项目都可以通过基本
**   项目的闭包来获得。
**
** 计算Follow Set的3条规则
** 1. 如果start是开始符号，那么$在Follow(start)中
** 2. 如果有产生式 A -> pBq，那么 FIRST(q) 除了 Є 都在 Follow(A)中
** 3. 如果存在产生式 A->pB 或 A->pBq，其中FIRST(q) 包含 Є， 那么 FOLLOW(B) 或包含Follow(A)所有符号
*/

var freelist []*config  // List of free configurations
var current *config     // Top of list of configurations
var currentend **config // Last on list of configs TODO currentend指向的是current的后一个节点的指针？？ current所在列表的最后一个项目之后，即代表列表的末尾 ？？？
var basis *config       // Top of list of basis configs
var basisend **config   // End of list of basis configs TODO

// Return a pointer to a new configuration   TODO lemon.c中为了提高效率，每次都一次性分为3个
func newconfig() *config {
	var new *config
	if freelist == nil || len(freelist) == 0 { // TODO 其实是一个对象池的概念
		amt := 3
		freelist = make([]*config, 0)
		for i := 0; i < amt; i++ {
			freelist = append(freelist, &config{})
		}

		for i := 0; i < amt-1; i++ {
			freelist[i].next = freelist[i+1]
		}
	}
	new = freelist[0]
	freelist = freelist[1:]
	return new
}

// The configuration "old" is no longer use
func deleteconfig(old *config) { // TODO how to add to head of slice
	old.next = freelist[0]
	freelist = append([]*config{old}, freelist...)
}

func Configlist_init() {
	current = nil
	currentend = &current
	basis = nil
	basisend = &basis
	Configtable_init()
}

func Configlist_reset() {
	current = nil
	currentend = &current
	basis = nil
	basisend = &basis
	Configtable_clear(nil)
}

// Add another configuration to the configuration list
func Configlist_add(rp *rule, dot int) *config {
	var cfp *config
	var model *config

	//assert(currentend != 0);
	model.rp = rp
	model.dot = dot
	cfp = Configtable_find(model)
	if cfp == nil {
		cfp = newconfig()
		cfp.rp = rp
		cfp.dot = dot
		cfp.fws = SetNew()
		cfp.stp = nil
		cfp.fplp = nil
		cfp.bplp = nil
		cfp.next = nil
		cfp.bp = nil
		*currentend = cfp
		currentend = &cfp.next
		Configtable_insert(cfp)
	}
	return cfp
}

// Add a basis configuration to the configuration list
func Configlist_addbasis(rp *rule, dot int) *config {
	var cfp *config
	var model *config

	//assert(basisend != 0);
	//assert(currentend != 0);
	model.rp = rp
	model.dot = dot
	cfp = Configtable_find(model)
	if cfp == nil {
		cfp = newconfig()
		cfp.rp = rp
		cfp.dot = dot
		cfp.fws = SetNew()
		cfp.stp = nil
		cfp.fplp = nil
		cfp.bplp = nil
		cfp.next = nil
		cfp.bp = nil
		*currentend = cfp
		currentend = &cfp.next
		*basisend = cfp
		basisend = &cfp.bp
		Configtable_insert(cfp)
	}
	return cfp
}

// Compute the closure of the configuration list 计算闭包
func Configlist_closure(lemp *lemon) {
	var cfp, newcfp *config
	var rp, newrp *rule
	var sp, xsp *symbol
	var i, dot int

	// assert currentend != nil

	for cfp = current; cfp != nil; cfp = cfp.next {
		rp = cfp.rp
		dot = cfp.dot
		if dot >= rp.nrhs {
			continue
		}

		sp = rp.rhs[dot]
		if sp.typ == NONTERMINAL {
			if sp.rule == nil && sp != lemp.errsym {
				ErrorMsg(lemp.filename, rp.line,
					fmt.Sprintf("Nonterminal \"%s\" has no rules.", sp.name))
				lemp.errorcnt++;
			}
			for newrp = sp.rule; newcfp != nil; newrp = newrp.nextlhs {
				// Configlist_add会把newcfp加到以current为表头的链表的末尾，因此在循环的同时，链表也在增长，
				// 之后会在外围的for循环中访问到，因此当外围的循环结束时，current所在的链表上所有项目就组成了一个状态
				newcfp = Configlist_add(newrp, 0)

				// 更新newcfp的FollowSet， 计算产生式的第二条规则
				for i = dot + 1; i < rp.nrhs; i++ {
					xsp = rp.rhs[i]
					if xsp.typ == TERMINAL {
						SetAdd(newcfp.fws, xsp.index)
						break
					} else {
						SetUnion(newcfp.fws, xsp.firstset)
						if !xsp.lambda {
							break
						}
					}
				}

				// 计算产生式的第三条规则
				if i == rp.nrhs {
					Plink_add(&cfp.fplp, newcfp) // TODO
				}
			}
		}
	}
}

func Configlist_sort() {
	current = (*config)(msort(unsafe.Pointer(current), cmpCfg, getNextCfg, setNextCfg))
	currentend = nil
}

func Configlist_sortbasis() {
	basis = (*config)(msort(unsafe.Pointer(current), cmpCfg, getNextCfg, setNextCfg)) // TODO ???
	basisend = nil
}

// Return a pointer to the head of the configuration list and reset the list
func Configlist_return() *config {
	old := current
	current = nil
	currentend = nil
	return old
}

// Return a pointer to the head of the configuration list and reset the list
func Configlist_basis() *config {
	old := basis
	basis = nil
	basisend = nil
	return old
}

// Free all elements of the given configuration list
func Configlist_eat(cfp *config) {
	var nextcfp *config
	for ; cfp != nil; cfp = nextcfp {
		nextcfp = cfp.next
		//assert cfp.fplp == nil
		//assert cfp.bplp == nil
		if cfp.fws != nil && len(cfp.fws) > 0 {
			SetFree(cfp.fws)
			deleteconfig(cfp)
		}
	}
}
