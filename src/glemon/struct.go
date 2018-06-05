package glemon

// Token 类型
type symbol_type int

const (
	TERMINAL      symbol_type = iota
	NONTERMINAL
	MULTITERMINAL
)

// 结合性
type e_assoc int

const (
	LEFT  e_assoc = iota
	RIGHT
	NONE
	UNK
)

/**
表示符号
 */
type symbol struct {
	name         string      /* Name of the symbol */
	index        int         /* Index number for this symbol */
	typ          symbol_type /* Symbols are all either TERMINALS or NTs */
	rule         *rule       /* Linked list of rules of this (if an NT) 如果是非终结符，那么会有多个产生式，symbol的rule是链表的表首指针， rule结构体里的nextlhs指向下一个节点*/
	fallback     *symbol     /* fallback token in case this token doesn't parse */
	prec         int         /* Precedence if defined (-1 otherwise) */
	assoc        e_assoc     /* Associativity if precedence is defined */
	firstset     []byte      /* First-set for all rules of this symbol */
	lambda       bool        /* True if NT and can generate an empty string */
	destructor   string      /* Code which executes whenever this symbol is popped from the stack during error processing */
	destructorln int         /* Line number for start of destructor.  Set to -1 for duplicate destructors. */
	datatype     string      /* The data type of information held by this object. Only used if type==NONTERMINAL */
	dtnum        int         /* The data type number.  In the parser, the value stack is a union.  The .yy%d element of this union is the correct data type for this object */
}

/* Each production rule in the grammar is stored in the following structure.
** 表示产生式  */
type rule struct {
	lhs       *symbol   /* Left-hand side of the rule 产生式左边的非终结符*/
	lhsalias  string    /* Alias for the LHS (NULL if none) 产生式左边的非终结符的别名*/
	ruleline  int       /* Line number for the rule */
	nrhs      int       /* Number of RHS symbols */
	rhs       []*symbol /* The RHS symbols 产生式右边所有*/
	rhsalias  []string  /* An alias for each RHS symbol (NULL if none) 产生式右边的各符号的别名*/
	line      int       /* Line number at which code begins */
	code      string    /* The code executed when this rule is reduced 规约时执行的动作代码*/
	precsym   *symbol   /* Precedence symbol for this rule */
	index     int       /* An index number for this rule */
	canReduce bool      /* True if this rule is ever reduced */
	nextlhs   *rule     /* Next rule with the same LHS 具有相同左边非终结符的下个产生式*/
	next      *rule     /* Next rule in the global list */
}

/* A configuration is a production rule of the grammar together with
** a mark (dot) showing how much of that rule has been processed so far.
** Configurations also contain a follow-set which is a list of terminal
** symbols which are allowed to immediately follow the end of the rule.
** Every configuration is recorded as an instance of the following:
*/
type cfgstatus int

const (
	COMPLETE   cfgstatus = iota
	INCOMPLETE
)

/*
** config用于保存右部加有句点的产生式，即LALR中的item */
type config struct {
	rp     *rule     /* The rule upon which the configuration is based */
	dot    int       /* The parse point */
	fws    []byte    /* Follow-set for this configuration only */
	fplp   *plink    /* Follow-set forward propagation links */
	bplp   *plink    /* Follow-set backwards propagation links */
	stp    *state    /* Pointer to state which contains this */
	status cfgstatus /* used during followset and shift computations */
	next   *config   /* Next configuration in the state */
	bp     *config   /* The next basis configuration */
}

type e_action int

const (
	SHIFT       e_action = iota
	ACCEPT
	REDUCE
	ERROR
	CONFLICT
	SH_RESOLVED  /* Was a shift.  Precedence resolved conflict */
	RD_RESOLVED  /* Was reduce.  Precedence resolved conflict */
	NOT_USED     /* Deleted by compression */
)

/* Every shift or reduce operation is stored as one of the following */
type action struct {
	sp      *symbol /* The look-ahead symbol */
	typ     e_action
	stp     *state  /* The new state, if a shift */
	rp      *rule   /* The rule, if a reduce */
	next    *action /* Next action for this state */
	collide *action /* Next action with the same hash */
}

/* Each state of the generated parser's finite state machine
** is encoded as an instance of the following structure. */
type state struct {
	bp                *config /* The basis configurations for this state */
	cfp               *config /* All configurations in this set */
	index             int     /* Sequential number for this state */
	ap                *action /* List of actions for this state */
	nTknAct, nNtAct   int     /* Number of actions on terminals and nonterminals */
	iTknOfst, iNtOfst int     /* yy_action[] offset for terminals and nonterminals */
	iDflt             int     /* Default action */
}

/* A followset propagation link indicates that the contents of one
** configuration followset should be propagated to another whenever
** the first changes. */
type plink struct {
	cfp  *config /* The configuration to which linked */
	next *plink  /* The next propagate link */
}

//The state vector for the entire parser generator is recorded asfollows.
type lemon struct {
	sorted []*state /* Table of states sorted by state number */
	rule   *rule    /* List of all rules rule组成一个链表，此处是首节点 */

	nstate       int       /* Number of states */
	nrule        int       /* Number of rules */
	nsymbol      int       /* Number of terminal and nonterminal symbols */
	nterminal    int       /* Number of terminal symbols */
	symbols      []*symbol /* Sorted array of pointers to symbols */
	errorcnt     int       /* Number of errors */
	errsym       *symbol   /* The error symbol */
	name         string    /* Name of the generated parser */
	arg          string    /* Declaration of the 3th argument to parser */
	tokentype    string    /* Type of terminal symbols in the parser stack */
	vartype      string    /* The default type of non-terminal symbols */
	start        string    /* Name of the start symbol for the grammar */
	stacksize    string    /* Size of the parser stack (存储为string, 是为了parseonetoken()方便处理)*/
	include      string    /* Code to put at the start of the C file */
	includeln    int       /* Line number for start of include code */
	error        string    /* Code to execute when an error is seen */
	errorln      int       /* Line number for start of error code */
	overflow     string    /* Code to execute on a stack overflow */
	overflowln   int       /* Line number for start of overflow code */
	failure      string    /* Code to execute on parser failure */
	failureln    int       /* Line number for start of failure code */
	accept       string    /* Code to execute when the parser excepts */
	acceptln     int       /* Line number for the start of accept code */
	extracode    string    /* Code appended to the generated file */
	extracodeln  int       /* Line number for the start of the extra code */
	tokendest    string    /* Code to execute to destroy token data */
	tokendestln  int       /* Line number for token destroyer code */
	vardest      string    /* Code for the default non-terminal destructor */
	vardestln    int       /* Line number for default non-term destructor code*/
	filename     string    /* Name of the input file */
	outname      string    /* Name of the current output file */
	tokenprefix  string    /* A prefix added to token names in the .h file */
	nconflict    int       /* Number of parsing conflicts */
	tablesize    int       /* Size of the parse tables */
	basisflag    int       /* Print only basis configurations */
	has_fallback bool      /* True if any %fallback is seen in the grammer */
	argv0        string    /* Name of the program */
}
