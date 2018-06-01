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

// 解析语法定义文件的内容为token流
func Parse(gp *lemon) {
	ps := pstate{}
	ps.gp = gp
	ps.filename = gp.filename
	ps.state = INITIALIZE
	startline := 0

	filebuf, err := ioutil.ReadFile(ps.filename)
	if err != nil {
		ErrorMsg(ps.filename, 0, "Can't read file.");
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
					"String starting on this line is not terminated before the end of the file.");
				ps.errorcnt++;
				nextcp = cp;
			} else {
				nextcp = cp + 1;
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
					"C code starting on this line is not terminated before the end of the file.");
				ps.errorcnt++;
				nextcp = cp;
			} else {
				nextcp = cp + 1;
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

		// TODO 按照这个逻辑对于字符串， "abcd", token的内容会变成 "abcd , 会丢失最后面的双引号，对于代码块会丢失最好的右括弧
	}

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
		psp.lastrule = nilp
		psp.gp.nrule = 0
		fallthrough
	case WAITING_FOR_DECL_OR_RULE:
		if x[0] == '%' { // 特殊申明
			psp.state = WAITING_FOR_DECL_KEYWORD
		} else if util.IsLowerChar(x[0]) { //产生左边
			psp.lhs = Symbol_new(token)
			psp.nrhs = 0
			psp.lhsalias = ""
			psp.state = WAITING_FOR_ARROW
		} else if x[0] == '{' { // 动作代码
			if psp.prevrule == nil {
				ErrorMsg(psp.filename, psp.tokenlineno, "There is not prior rule "+
					"upon which to attach the code fragment which begins on this line.");
				psp.errorcnt++;
			} else if len(psp.prevrule.code) > 0 { // 已经指定过动作代码了
				ErrorMsg(psp.filename, psp.tokenlineno, "Code fragment beginning on this line "+
					"is not the first to follow the previous rule.");
				psp.errorcnt++;
			} else {
				psp.prevrule.line = psp.tokenlineno
				psp.prevrule.code = string(x[1:]) // 不需要最外面的括弧
			}
		} else if x[0] == '[' {
			psp.state = PRECEDENCE_MARK_1 // 产生式后边用于指定优先级的符号
		} else {
			ErrorMsg(psp.filename, psp.tokenlineno,
				fmt.Sprintf("Token \"%s\" should be either \"%%\" or a nonterminal name.", x));
			psp.errorcnt++;
		}
	case PRECEDENCE_MARK_1:
		if !util.IsUpperChar(x[0]) {
			ErrorMsg(psp.filename, psp.tokenlineno,
				"The precedence symbol must be a terminal.");
			psp.errorcnt++;
		} else if psp.prevrule == nil {
			ErrorMsg(psp.filename, psp.tokenlineno,
				"There is no prior rule to assign precedence \"[%s]\".", x);
			psp.errorcnt++;
		} else if psp.prevrule.precsym != nil {
			ErrorMsg(psp.filename, psp.tokenlineno, "Precedence mark on this line "+
				"is not the first to follow the previous rule.");
			psp.errorcnt++;
		} else {
			psp.prevrule.precsym = Symbol_new(token)
		}
		psp.state = PRECEDENCE_MARK_2

		// TODO Line 1964
	}

	fmt.Printf("%s, %d \n", ps.tokenstart, ps.tokenlineno)
	fmt.Printf("%s\n", "------------------------")
}
