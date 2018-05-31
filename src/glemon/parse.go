package glemon

/*
** Input file parser for the LEMON parser generator.
** 词法分析器
*/

type e_state int

const (
	INITIALIZE                    e_state = iota
	WAITING_FOR_DECL_OR_RULE
	WAITING_FOR_DECL_KEYWORD
	WAITING_FOR_DECL_ARG
	WAITING_FOR_PRECEDENCE_SYMBOL
	WAITING_FOR_ARROW
	IN_RHS
	LHS_ALIAS_1
	LHS_ALIAS_2
	LHS_ALIAS_3
	RHS_ALIAS_1
	RHS_ALIAS_2
	PRECEDENCE_MARK_1
	PRECEDENCE_MARK_2
	RESYNC_AFTER_RULE_ERROR
	RESYNC_AFTER_DECL_ERROR
	WAITING_FOR_DESTRUCTOR_SYMBOL
	WAITING_FOR_DATATYPE_SYMBOL
	WAITING_FOR_FALLBACK_ID
)

type pstate struct {
	filename    string    /* Name of the input file */
	tokenlineno int       /* Linenumber at which current token starts */
	errorcnt    int       /* Number of errors so far */
	tokenstart  string    /* Text of current token 字符串形式的符号名称*/
	gp          *lemon    /* Global state vector */
	state       e_state   /* The state of the parser */
	fallback    *symbol   /* The fallback token */
	lhs         *symbol   /* Left-hand side of current rule */
	lhsalias    string    /* Alias for the LHS */
	nrhs        int       /* Number of right-hand side symbols seen */
	rhs         []*symbol /* RHS symbols */
	alias       []string  /* Aliases for each RHS symbol (or NULL) */
	prevrule    *rule     /* Previous rule parsed */
	declkeyword string    /* Keyword of a declaration */
	declargslot []string  /* Where the declaration argument should be put TODO */
	decllnslot  []int     /* Where the declaration linenumber is put */
	declassoc   e_assoc   /* Assign this association to decl arguments */
	preccounter int       /* Assign this precedence to decl arguments */
	firstrule   *rule     /* Pointer to first rule in the grammar */
	lastrule    *rule     /* Pointer to the most recently parsed rule */
}

















