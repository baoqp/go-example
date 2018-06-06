package glemon

import (
	"io/ioutil"
	"unicode"
	"fmt"
	"util"
)

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

const (
	MAXRHS = 1000
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
	declargslot *string   /* Where the declaration argument should be put TODO */
	decllnslot  *int      /* Where the declaration linenumber is put */
	declassoc   e_assoc   /* Assign this association to decl arguments */
	preccounter int       /* Assign this precedence to decl arguments */
	firstrule   *rule     /* Pointer to first rule in the grammar */
	lastrule    *rule     /* Pointer to the most recently parsed rule */
}

// 解析语法定义文件的内容为token流
func Parse(gp *lemon) {
	ps := pstate{}
	ps.gp = gp
	ps.filename = gp.filename
	ps.state = INITIALIZE
	startline := 0

	filebuf, err := ioutil.ReadFile(ps.filename)
	if err != nil {
		ErrorMsg(ps.filename, 0, "Can't read file.")
		gp.errorcnt ++
		return
	}

	//Now scan the text of the input file
	lineno := 1
	bufLen := len(filebuf)
	var cp, nextcp int
	for cp = 0; cp < bufLen; {

		if filebuf[cp] == '\n' {
			lineno ++ //Keep track of the line number
		}
		if unicode.IsSpace(rune(filebuf[cp])) {
			cp ++
			continue
		}

		if filebuf[cp] == '/' && filebuf[cp+1] == '/' { // Skip C style comments
			cp += 2
			for ; cp < bufLen && filebuf[cp] != '\n'; {
				cp++
			}

			continue
		}

		if filebuf[cp] == '/' && filebuf[cp+1] == '*' { // Skip C style comments
			cp += 2
			for ; cp < bufLen && ( filebuf[cp] != '/' || filebuf[cp-1] != '*' ); {
				cp++

				if filebuf[cp] == '\n' {
					lineno ++
				}
			}
			if cp < bufLen {
				cp++
			}
			continue
		}

		startIdx := cp // TODO

		ps.tokenlineno = lineno // Linenumber on which token begins

		if filebuf[cp] == '"' { // String literals
			cp ++
			for ; cp < bufLen && filebuf[cp] != '"'; {
				cp++
				if filebuf[cp] == '\n' {
					lineno ++
				}
			}
			if cp == bufLen {
				ErrorMsg(ps.filename, startline,
					"String starting on this line is not terminated before the end of the file.")
				ps.errorcnt++
				nextcp = cp
			} else {
				nextcp = cp + 1
			}
		} else if filebuf[cp] == '{' { // code block
			var level int // level 用于表示有几层{}
			cp ++
			for level = 1; cp < bufLen && (level > 1 || filebuf[cp] != '}'); cp ++ {
				if filebuf[cp] == '\n' {
					lineno ++
				} else if filebuf[cp] == '{' {
					level ++
				} else if filebuf[cp] == '}' {
					level --
				} else if filebuf[cp] == '/' && filebuf[cp+1] == '*' { // Skip comments
					cp += 2
					var prevc byte = 0
					for ; cp < bufLen && ( filebuf[cp] != '/' || prevc != '*' ); {
						cp++
						prevc = filebuf[cp]
						if filebuf[cp] == '\n' {
							lineno ++
						}
					}
				} else if filebuf[cp] == '/' && filebuf[cp+1] == '/' { // Skip comments
					cp += 2
					for ; cp < bufLen && ( filebuf[cp] != '\n' ); {
						cp ++
					}

					if cp < bufLen {
						lineno ++
					}
				} else if filebuf[cp] == '"' || filebuf[cp] == '\'' { // String a character literals
					startchar := filebuf[cp]
					prevc := byte(0)

					for cp++; cp < bufLen && (filebuf[cp] != startchar || prevc == '\\'); cp++ {
						if filebuf[cp] == '\n' {
							lineno++
						}
						if prevc == '\\' {
							prevc = 0
						} else {
							prevc = filebuf[cp]
						}
					}

				}

			}

			if cp == bufLen {
				ErrorMsg(ps.filename, ps.tokenlineno,
					"C code starting on this line is not terminated before the end of the file.")
				ps.errorcnt++
				nextcp = cp
			} else {
				nextcp = cp + 1
			}

		} else if util.IsAlumn(filebuf[cp]) { // Identifiers
			for ; cp < bufLen && (util.IsAlumn(filebuf[cp]) || filebuf[cp] == '_' ); {
				cp ++
			}
			nextcp = cp
		} else if filebuf[cp] == ':' && filebuf[cp+1] == ':' && filebuf[cp+2] == '=' {
			cp += 3
			nextcp = cp
		} else { // All other (one character) operators
			cp ++
			nextcp = cp
		}
		ps.tokenstart = string(filebuf[startIdx:cp])
		parseonetoken(&ps)
		cp = nextcp

		// 按照这个逻辑对于字符串， "abcd", token的内容会变成 "abcd , 会丢失最后面的双引号，对于代码块会丢失最好的右括弧
	}
	gp.rule = ps.firstrule
	gp.errorcnt = ps.errorcnt

}

//  Parse a single token
func parseonetoken(psp *pstate) {
	token := psp.tokenstart
	Strsafe(token)
	x := []byte(token)

	switch psp.state {
	case INITIALIZE:
		psp.prevrule = nil
		psp.preccounter = 0
		psp.firstrule = nil
		psp.lastrule = nil
		psp.gp.nrule = 0
		fallthrough
	case WAITING_FOR_DECL_OR_RULE:
		if x[0] == '%' { // 特殊申明
			psp.state = WAITING_FOR_DECL_KEYWORD
		} else if util.IsLowerChar(x[0]) { //产生左边
			psp.lhs = Symbol_new(token)
			psp.nrhs = 0
			psp.rhs = make([]*symbol, 0)
			psp.lhsalias = ""
			psp.alias = make([]string, 0)
			psp.state = WAITING_FOR_ARROW
		} else if x[0] == '{' { // 动作代码
			if psp.prevrule == nil {
				ErrorMsg(psp.filename, psp.tokenlineno, "There is not prior rule "+
					"upon which to attach the code fragment which begins on this line.")
				psp.errorcnt++
			} else if len(psp.prevrule.code) > 0 { // 已经指定过动作代码了
				ErrorMsg(psp.filename, psp.tokenlineno, "Code fragment beginning on this line "+
					"is not the first to follow the previous rule.")
				psp.errorcnt++
			} else {
				psp.prevrule.line = psp.tokenlineno
				psp.prevrule.code = string(x[1:]) // 不需要最外面的括弧
			}
		} else if x[0] == '[' {
			psp.state = PRECEDENCE_MARK_1 // 产生式后边用于指定优先级的符号
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Token \"%s\" should be either \"%%\" or a nonterminal name.", token))
			psp.errorcnt++
		}
	case PRECEDENCE_MARK_1:
		if !util.IsUpperChar(x[0]) {
			ErrorMsg(psp.filename, psp.tokenlineno,
				"The precedence symbol must be a terminal.")
			psp.errorcnt++
		} else if psp.prevrule == nil {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("There is no prior rule to assign precedence \"[%s]\".", token))
			psp.errorcnt++
		} else if psp.prevrule.precsym != nil {
			ErrorMsg(psp.filename, psp.tokenlineno, "Precedence mark on this line "+
				"is not the first to follow the previous rule.")
			psp.errorcnt++
		} else {
			psp.prevrule.precsym = Symbol_new(token)
		}
		psp.state = PRECEDENCE_MARK_2

	case PRECEDENCE_MARK_2:
		if x[0] != ']' {
			ErrorMsg(psp.filename, psp.tokenlineno,
				"Missing \"]\" on precedence mark.")
			psp.errorcnt++
		}
		psp.state = WAITING_FOR_DECL_OR_RULE
	case WAITING_FOR_ARROW:
		if x[0] == ':' && x[1] == ':' && x[2] == '=' {
			psp.state = IN_RHS
		} else if x[0] == '(' { // 别名
			psp.state = LHS_ALIAS_1
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Expected to see a \":\" following the LHS symbol \"%s\".", psp.lhs.name))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case LHS_ALIAS_1:
		if util.IsAlphaChar(x[0]) {
			psp.lhsalias = token
			psp.state = LHS_ALIAS_2
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("\"%s\" is not a valid alias for the LHS \"%s\"\n", x, psp.lhs.name))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case LHS_ALIAS_2:
		if x[0] == ')' {
			psp.state = LHS_ALIAS_3
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Missing \")\" following LHS alias name \"%s\".", psp.lhsalias))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case LHS_ALIAS_3:
		if x[0] == ':' && x[1] == ':' && x[2] == '=' {
			psp.state = IN_RHS
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Missing \"->\" following: \"%s(%s)\".", psp.lhs.name, psp.lhsalias))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case IN_RHS: // 进入产生式右边
		if x[0] == '.' { // 遇到句点，说明已经到了产生式的几位
			rp := &rule{}
			for i := 0; i < psp.nrhs; i++ {
				rp.rhs = append(rp.rhs, psp.rhs[i])
				rp.rhsalias = append(rp.rhsalias, psp.alias[i])
			}
			rp.lhs = psp.lhs
			rp.lhsalias = psp.lhsalias
			rp.nrhs = psp.nrhs
			rp.code = ""
			rp.precsym = nil
			rp.index = psp.gp.nrule
			psp.gp.nrule ++
			rp.nextlhs = rp.lhs.rule
			rp.lhs.rule = rp
			rp.next = nil
			if psp.firstrule == nil {
				psp.firstrule = rp
				psp.lastrule = rp
			} else {
				psp.lastrule.next = rp
				psp.lastrule = rp
			}
			psp.prevrule = rp
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if util.IsAlphaChar(x[0]) {
			if psp.nrhs >= MAXRHS {
				ErrorMsg(psp.filename, psp.tokenlineno,
					fmt.Sprintf("Too many symbol on RHS or rule beginning at \"%s\".", token))
				psp.errorcnt++
				psp.state = RESYNC_AFTER_RULE_ERROR
			} else {
				psp.rhs = append(psp.rhs, Symbol_new(token))
				psp.alias = append(psp.alias, "")
				psp.nrhs ++
			}
		} else if x[0] == '(' && psp.nrhs > 0 {
			psp.state = RHS_ALIAS_1
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Illegal character on RHS of rule: \"%s\".", token))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case RHS_ALIAS_1:
		if util.IsAlphaChar(x[0]) {
			psp.alias[psp.nrhs-1] = token
			psp.state = RHS_ALIAS_2
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("\"%s\" is not a valid alias for the RHS symbol \"%s\"\n",
					token, psp.rhs[psp.nrhs-1].name))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case RHS_ALIAS_2:
		if x[0] == ')' {
			psp.state = IN_RHS
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Missing \")\" following LHS alias name \"%s\".", psp.lhsalias))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
	case WAITING_FOR_DECL_KEYWORD:
		if util.IsAlphaChar(x[0]) {
			psp.declkeyword = token
			psp.declargslot = nil
			psp.decllnslot = nil
			psp.state = WAITING_FOR_DECL_ARG
			if token == "name" {
				psp.declargslot = &psp.gp.name
			} else if token == "include" {
				psp.declargslot = &psp.gp.include
				psp.decllnslot = &psp.gp.includeln
			} else if token == "code" {
				psp.declargslot = &psp.gp.extracode
				psp.decllnslot = &psp.gp.extracodeln
			} else if token == "token_prefix" {
				psp.declargslot = &psp.gp.tokenprefix
			} else if token == "syntax_error" {
				psp.declargslot = &psp.gp.error
				psp.decllnslot = &psp.gp.errorln
			} else if token == "parse_accept" {
				psp.declargslot = &psp.gp.accept
				psp.decllnslot = &psp.gp.acceptln
			} else if token == "parse_failure" {
				psp.declargslot = &psp.gp.failure
				psp.decllnslot = &psp.gp.failureln
			} else if token == "stack_overflow" {
				psp.declargslot = &psp.gp.overflow
				psp.decllnslot = &psp.gp.overflowln
			} else if token == "extra_argument" {
				psp.declargslot = &psp.gp.arg
			} else if token == "token_type" {
				psp.declargslot = &psp.gp.tokentype
			} else if token == "default_type" {
				psp.declargslot = &psp.gp.vartype
			} else if token == "stack_size" {
				psp.declargslot = &psp.gp.stacksize
			} else if token == "start_symbol" {
				psp.declargslot = &psp.gp.start
			} else if token == "left" {
				psp.preccounter++
				psp.declassoc = LEFT
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			} else if token == "right" {
				psp.preccounter++
				psp.declassoc = RIGHT
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			} else if token == "nonassoc" {
				psp.preccounter++
				psp.declassoc = NONE
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			} else if token == "type" {
				psp.state = WAITING_FOR_DATATYPE_SYMBOL
			} else if token == "fallback" {
				psp.fallback = nil
				psp.state = WAITING_FOR_FALLBACK_ID
			} else { // do not support token_destructor, default_destructor or destructor
				ErrorMsg(psp.filename, psp.tokenlineno,
					fmt.Sprintf("Unknown declaration keyword: \"%%%s\".", token))
				psp.errorcnt++
				psp.state = RESYNC_AFTER_DECL_ERROR
			}

		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Illegal declaration keyword: \"%s\".", token))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		}
	case WAITING_FOR_DATATYPE_SYMBOL:
		if !util.IsAlphaChar(x[0]) {
			ErrorMsg(psp.filename, psp.tokenlineno,
				"Symbol name missing after %type keyword")
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		} else {
			sp := Symbol_new(token)
			psp.declargslot = &sp.datatype
			psp.decllnslot = nil
			psp.state = WAITING_FOR_DECL_ARG
		}
	case WAITING_FOR_PRECEDENCE_SYMBOL:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if util.IsUpperChar(x[0]) {
			sp := Symbol_new(token)
			if sp.prec >= 0 {
				ErrorMsg(psp.filename, psp.tokenlineno,
					fmt.Sprintf("Symbol \"%s\" has already be given a precedence.", token))
				psp.errorcnt++
			} else {
				sp.prec = psp.preccounter
				sp.assoc = psp.declassoc
			}
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Can't assign a precedence to \"%s\".", token))
			psp.errorcnt++
		}
	case WAITING_FOR_DECL_ARG:
		if x[0] == '{' || x[0] == '"' || util.IsAlumn(x[0]) {
			if  psp.declargslot != nil && len(*psp.declargslot) > 0 {
				errmsg := token
				if x[0] == '"' {
					errmsg = token[1:]
				}
				ErrorMsg(psp.filename, psp.tokenlineno,
					fmt.Sprintf("The argument \"%s\" to declaration \"%%%s\" is not the first.",
						errmsg, psp.declkeyword))
				psp.errorcnt++
				psp.state = RESYNC_AFTER_DECL_ERROR
			} else {

				*psp.declargslot = token
				if x[0] == '"' || x[0] == '{' {
					*(psp.declargslot) = token[1:]
				}

				if psp.decllnslot != nil {
					*psp.decllnslot = psp.tokenlineno
				}
				psp.state = WAITING_FOR_DECL_OR_RULE
			}
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Illegal argument to %%%s: %s", psp.declkeyword, token))
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		}

	case WAITING_FOR_FALLBACK_ID:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if !util.IsUpperChar(x[0]) {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("%%fallback argument \"%s\" should be a token", token))
			psp.errorcnt++
		} else {
			sp := Symbol_new(token)
			if psp.fallback == nil {
				psp.fallback = sp
			} else if sp.fallback != nil {
				ErrorMsg(psp.filename, psp.tokenlineno,
					fmt.Sprintf("More than one fallback assigned to token %s", token))
				psp.errorcnt++
			} else {
				sp.fallback = psp.fallback
				psp.gp.has_fallback = true
			}
		}
	case RESYNC_AFTER_RULE_ERROR:
	case RESYNC_AFTER_DECL_ERROR:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		}
		if x[0] == '%' {
			psp.state = WAITING_FOR_DECL_KEYWORD
		}
	}

	//fmt.Printf("%s, %d \n", psp.tokenstart, psp.tokenlineno)
	//fmt.Printf("%s \n", "------------------------")
}
