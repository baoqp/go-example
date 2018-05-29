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
	name       string      /* Name of the symbol */
	index      int         /* Index number for this symbol */
	typ        symbol_type /* Symbols are all either TERMINALS or NTs */
	rule       *rule       /* Linked list of rules of this (if an NT) */
	fallback   *symbol     /* fallback token in case this token doesn't parse */
	prec       int         /* Precedence if defined (-1 otherwise) */
	assoc      e_assoc     /* Associativity if precedence is defined */
	firstset   []rune      /* First-set for all rules of this symbol */
	lambda     bool        /* True if NT and can generate an empty string */
	useCnt     int         /* Number of times used */
	destructor string      /* Code which executes whenever this symbol is popped from the stack during error processing */
	destLineno int         /* Line number for start of destructor.  Set to -1 for duplicate destructors. */
	datatype   string      /* The data type of information held by this object. Only used if type==NONTERMINAL */
	dtnum      int         /* The data type number.  In the parser, the value stack is a union.  The .yy%d element of this union is the correct data type for this object */
	bContent   int         /* True if this symbol ever carries content - if it is ever more than just syntax */
	/* The following fields are used by MULTITERMINALs only */
	nsubsym int       /* Number of constituent symbols in the MULTI */
	subsym  []*symbol /* Array of constituent symbols */
}

/* Each production rule in the grammar is stored in the following structure.
** 表示产生式  */
type rule struct {
	lhs         *symbol   /* Left-hand side of the rule */
	lhsalias    string    /* Alias for the LHS (NULL if none) */
	lhsStart    bool      /* True if left-hand side is the start symbol */
	ruleline    int       /* Line number for the rule */
	nrhs        int       /* Number of RHS symbols */
	rhs         []*symbol /* The RHS symbols */
	rhsalias    []string; /* An alias for each RHS symbol (NULL if none) */
	line        int       /* Line number at which code begins */
	code        string    /* The code executed when this rule is reduced */
	codePrefix  string    /* Setup code before code[] above */
	codeSuffix  string    /* Breakdown code after code[] above */
	noCode      bool      /* True if this rule has no associated C code */
	codeEmitted int       /* True if the code has been emitted already */
	precsym     *symbol   /* Precedence symbol for this rule */
	index       int       /* An index number for this rule */
	iRule       int       /* Rule number as used in the generated tables */
	canReduce   bool      /* True if this rule is ever reduced */
	doesReduce  bool      /* Reduce actions occur after optimization */
	nextlhs     *rule     /* Next rule with the same LHS */
	next        rule      /* Next rule in the global list */
}

//The state vector for the entire parser generator is recorded asfollows.
type lemon struct {
}

/* Each state of the generated parser's finite state machine
** is encoded as an instance of the following structure. */
type state struct {
}
