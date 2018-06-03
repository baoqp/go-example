package glemon

/*
** Routines to processing a configuration list and building a state
** in the LEMON parser generator.
** config是其他语法书中的LR item
*/

var freelist []config    // List of free configurations
var current *config     // Top of list of configurations
var currentend **config // Last on list of configs TODO
var basis *config       // Top of list of basis configs
var basisend **config   // End of list of basis configs TODO

// Return a pointer to a new configuration   TODO lemon.c中为了提高效率，每次都一次性分为3个
func newconfig() *config {
	var new *config
	if freelist == nil || len(freelist) == 0 { // TODO 其实是一个对象池的概念
		amt := 3
		freelist = make([]config,0)
		for i:=0; i<amt; i++ {
			freelist = append(freelist, config{})
		}

		for i:=0; i < amt -1; i++ {
			freelist[i].next = &freelist[i + 1]
		}
	}
	new = &freelist[0]
	freelist = freelist[1:]
	return new
}

// The configuration "old" is no longer use
func deleteconfig(old *config) { // TODO how to add to head of slice
	freelist = append(freelist, *old)
}


func Configlist_init() {
	current = nil
	currentend = &current
	basis = nil
	basisend = &basis
	Configtable_init()
}

// Add a basis configuration to the configuration list
func Configlist_addbasis(rp *rule, dot int) {

	var cfp *config
	var model *config
	//assert(basisend != 0);
	//assert(currentend != 0);

	model.rp = rp
	model.dot = dot
	cfp = Configtable_find(model)
	if cfp == nil {

	}

}
