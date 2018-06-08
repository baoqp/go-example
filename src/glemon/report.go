package glemon

import (
	"os"
	"strings"
	"fmt"
	"util"
	"bufio"
	"io"

	"strconv"
)

/*
** Procedures for generating reports and tables in the LEMON parser generator.
**
** Generate a filename with the given suffix.  Space to hold the
** name comes from malloc() and must be freed by the calling
** function.
*/

const (
	LINESIZE = 1000
)

func file_makename(lemp *lemon, suffix string) string {
	name := lemp.filename
	cp := strings.LastIndexByte(name, '.')
	if cp != -1 {
		name = name[:cp]
	}
	name = name + suffix
	return name
}

func file_open(lemp *lemon, suffix string, flag int) *os.File {
	lemp.outname = file_makename(lemp, suffix)
	fp, err := os.OpenFile(lemp.outname, flag, 0666)
	if err != nil {
		fmt.Printf("cannot open file %s \n", lemp.outname)
		return nil
	}
	return fp
}

func ConfigPrint(fp *os.File, cfp *config) {
	rp := cfp.rp
	fmt.Fprintf(fp, "%s ::=", rp.lhs.name)
	for i := 0; i <= rp.nrhs; i++ {
		if i == cfp.dot {
			fmt.Fprintf(fp, " *")
		}
		if i == rp.nrhs {
			break;
		}
		fmt.Fprintf(fp, " %s", rp.rhs[i].name)
	}
}

func SetPrint(out *os.File, set []byte, lemp *lemon) {
	spacer := ""
	fmt.Fprintf(out, "%12s[", "")
	for i := 0; i < lemp.nterminal; i++ {
		if set[i] > 0 {
			fmt.Fprintf(out, "%s%s", spacer, lemp.symbols[i].name)
			spacer = " ";
		}
	}
	fmt.Fprintf(out, "]\n");
}

func PlinkPrint(out *os.File, plp *plink, tag string) {
	for plp != nil {
		fmt.Fprintf(out, "%12s%s (state %2d) ", "", tag, plp.cfp.stp.index)
		ConfigPrint(out, plp.cfp)
		fmt.Fprintf(out, "\n")
		plp = plp.next
	}
}

func PrintAction(ap *action, fp *os.File, indent int) bool {
	result := true
	switch ap.typ {
	case SHIFT:
		fmt.Fprintf(fp, "%*s shift   %d", indent, ap.sp.name, ap.stp.index)
	case REDUCE:
		fmt.Fprintf(fp, "%*s reduce  %d", indent, ap.sp.name, ap.rp.index)
	case ACCEPT:
		fmt.Fprintf(fp, "%*s accept", indent, ap.sp.name)
	case ERROR:
		fmt.Fprintf(fp, "%*s error", indent, ap.sp.name)
	case CONFLICT:
		fmt.Fprintf(fp, "%*s reduce %-3d ** Parsing conflict **",
			indent, ap.sp.name, ap.rp.index)
	case SH_RESOLVED:
	case RD_RESOLVED:
	case NOT_USED:
		result = false
	}

	return result
}

// Generate the "y.output" log file 打印所有的状态
func ReportOutput(lemp *lemon) {
	var stp *state
	var cfp *config
	var ap *action
	var fp *os.File

	fp = file_open(lemp, ".out", os.O_WRONLY|os.O_CREATE)
	if fp == nil {
		return
	}
	fmt.Fprintf(fp, " \b")

	for i := 0; i < lemp.nstate; i++ {
		stp = lemp.sorted[i]

		fmt.Fprintf(fp, "State %d:\n", stp.index)

		if lemp.basisflag > 0 {
			cfp = stp.bp
		} else {
			cfp = stp.cfp
		}

		for cfp != nil {

			if cfp.dot == cfp.rp.nrhs {
				buf := fmt.Sprintf("(%d)", cfp.rp.index)
				fmt.Fprintf(fp, "    %5s ", buf)
			} else {
				fmt.Fprintf(fp, "          ")
			}

			ConfigPrint(fp, cfp)
			fmt.Fprintf(fp, "\n")

			if TEST {
				SetPrint(fp, cfp.fws, lemp)
				PlinkPrint(fp, cfp.fplp, "TO ")
				PlinkPrint(fp, cfp.bplp, "FROM ")
			}

			if lemp.basisflag > 0 {
				cfp = cfp.bp
			} else {
				cfp = cfp.next
			}

		}

		fmt.Fprintf(fp, "\n")
		for ap = stp.ap; ap != nil; ap = ap.next {
			if PrintAction(ap, fp, 30) {
				fmt.Fprintf(fp, "\n")
			}
		}
		fmt.Fprintf(fp, "\n")
	}
	fp.Close()
}

// Generate a header file for the parser
func ReportHeader(lemp *lemon) {

	var prefix string

	if len(lemp.tokenprefix) > 0 {
		prefix = lemp.tokenprefix
	} else {
		prefix = ""
	}

	out := file_open(lemp, "_def.go", os.O_WRONLY|os.O_APPEND|os.O_CREATE)
	if out != nil {
		for i := 1; i < lemp.nterminal; i++ {
			fmt.Fprintf(out, "const %s%-30s = %2d\n", prefix, lemp.symbols[i].name, i)

		}
	}
}

// Reduce the size of the action tables, if possible, by making use of defaults.
//
// In this version, we take the most frequent REDUCE action and make
// it the default.  Only default a reduce if there are more than one.
// 把每个状态最频繁出现的归约动作作为该状态的默认（default）动作
func CompressTables(lemp *lemon) {
	fmt.Println("---compress table---")
	var stp *state
	var ap, ap2 *action
	var rp, rp2, rbest *rule
	var n, nbest int

	for i := 0; i < lemp.nstate; i++ {
		stp = lemp.sorted[i]
		nbest = 0
		rbest = nil

		for ap = stp.ap; ap != nil; ap = ap.next {
			if ap.typ != REDUCE {
				continue
			}

			rp = ap.rp
			if rp == rbest {
				continue
			}

			n = 1
			for ap2 = ap.next; ap2 != nil; ap2 = ap2.next {
				if ap2.typ != REDUCE {
					continue
				}

				rp2 = ap2.rp

				if rp2 == rbest {
					continue
				}
				if rp2 == rp {
					n++
				}

				if n > nbest {
					nbest = n
					rbest = rp
				}
			}
		}

		// Do not make a default if the number of rules to default is not at least 2
		if nbest < 2 {
			continue
		}

		for ap = stp.ap; ap != nil; ap = ap.next {
			if ap.typ == REDUCE && ap.rp == rbest {
				break
			}
		}

		// assert ap != nil
		// 和default相同的标记为不再使用，然后再对动作链表进行排序，由于{default}
		// 以'{'开头，其ASCII码值比如任何字母都大，因为会排在后面。在使用时，当遇到
		// 一个lookahead符号，会依次检查各个动作，都不符合最后会采用默认动作
		ap.sp = Symbol_new("{default}")
		for ap = ap.next; ap != nil; ap = ap.next {
			if ap.typ == REDUCE && ap.rp == rbest {
				ap.typ = NOT_USED
			}
		}

		stp.ap = Action_sort(stp.ap)

	}
}

func compute_action(lemp *lemon, ap *action) int {
	var act int
	switch ap.typ {
	case SHIFT:
		act = ap.stp.index
	case REDUCE:
		act = ap.rp.index + lemp.nstate
	case ERROR:
		act = lemp.nstate + lemp.nrule
	case ACCEPT:
		act = lemp.nstate + lemp.nrule + 1
	default:
		act = -1

	}
	return act
}

// transfers data from "in" to "out" until  a line is seen which begins with "%%".
// The line number is  tracked.  if name!=0, then any word that begin with "Parse"
// is changed to begin with name instead.
// 从模板文件中把内容转移到要生成的语法分析器文件中，直到遇到以%%开头的行。参数name由语法定义文件中的%name指定
var buff *bufio.Reader
func tplt_xfer(name string, in *os.File, out *os.File, lineno *int) {
	if buff == nil {
		buff  = bufio.NewReader(in) //读入缓存
	}
	for {
		line, err := buff.ReadString('\n') // 以'\n'为结束符读入一行
		if err != nil || io.EOF == err || (line[0] == '%' && line[1] == '%') {
			break
		}

		*lineno += 1
		iStart := 0

		if len(name) > 0 {
			for i := 0; i < len(line); i++ {
				// 把"Parse"替换成name
				if line[i] == 'P' && line[i:i+5] == "Parse" &&
					(i == 0 || !util.IsAlphaChar(line[i-1])) {
					if i > iStart {
						fmt.Fprintf(out, "%.*s", i-iStart, line[iStart:])
					}
					fmt.Fprintf(out, "%s", name);
					i += 4
					iStart = i + 1
				}
			}
		}
		fmt.Fprintf(out, "%s", line[iStart:])
	}
}

// finds the template file and opens it, returning a pointer to the opened file.
// 打开模板文件
func tplt_open(lemp *lemon) *os.File {
	templatename := "lempar.lt"
	var tpltname string

	// 如果存在自定义的模板文件，则使用它，否则使用默认的模板文件
	if idx := strings.LastIndexByte(lemp.filename, '.'); idx > 0 {
		tpltname = fmt.Sprintf("%s.lt", lemp.filename[:idx])
	} else {
		tpltname = fmt.Sprintf("%s.lt", lemp.filename)
	}

	if exists, err := util.Exists(tpltname); err != nil || !exists {
		tpltname = templatename
	}

	fp, err := os.OpenFile(tpltname, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Printf("Can't open the template file \"%s\".\n", templatename)
		lemp.errorcnt ++
		return nil
	}

	return fp
}

// Print a string to the file and keep the linenumber up to date
func tplt_print(out *os.File, lemp *lemon, str string, strln int, lineno *int) {
	if len(str) == 0 {
		return
	}

	fmt.Fprintf(out, "//#line %d \"%s\"\n", strln, lemp.filename);
	*lineno += 1

	for _, c := range str {
		fmt.Fprintf(out, "%c", c)
		if c == '\n' {
			*lineno += 1
		}
	}
	fmt.Fprintf(out, "\n//#line %d \"%s\"\n", *lineno+2, lemp.outname)
	*lineno += 2
}

// emits code for the destructor for the symbol sp TODO do not need destructor
func emit_destructor_code(out *os.File, sp *symbol, lemp *lemon, lineno *int) {

}

// Generate code which executes when the rule "rp" is reduced.  Write
// the code to "out".  Make sure lineno stays up-to-date.
func emit_code(out *os.File, rp *rule, lemp *lemon, lineno *int) {
	var cp byte
	var lhsused bool
	var linecnt int
	var used []bool = make([]bool, MAXRHS, MAXRHS)

	if len(rp.code) > 0 {
		fmt.Fprintf(out, "#line %d \"%s\"\n{", rp.line, lemp.filename)
		codes := []byte(rp.code)

		for i := 0; i < len(codes); i++ {
			cp = codes[i]
			if util.IsAlphaChar(cp) &&
				(i == 0 || (!util.IsAlumn(cp) && cp != '_')) {

				var j = i + 1
				for ; util.IsAlumn(rp.code[j]) || rp.code[j] == '_'; j++ {
				}

				if len(rp.lhsalias) > 0 && rp.code[i:j] == rp.lhsalias {
					fmt.Fprintf(out, "yygotominor.yy%d", rp.lhs.dtnum);
					i = j
					lhsused = true
				} else {
					for h := 0; h < rp.nrhs; h++ {
						if len(rp.rhsalias[h]) > 0 && rp.code[i:j] == rp.rhsalias[h] {
							fmt.Fprintf(out, "yymsp[%d].minor.yy%d", h-rp.nrhs+1, rp.rhs[h].dtnum);
							i = j
							used[h] = true
							break
						}
					}
				}
			}

			if cp == '\n' {
				linecnt ++
			}
			fmt.Fprintf(out, "%s", cp)
		}
		(*lineno) += 3 + linecnt
	} /* End if( rp.code ) */

	// Check to make sure the LHS has been used

	if len(rp.lhsalias) > 0 && !lhsused {
		ErrorMsg(lemp.filename, rp.ruleline,
			"Label \"%s\" for \"%s(%s)\" is never used.",
			rp.lhsalias, rp.lhs.name, rp.lhsalias)
		lemp.errorcnt++
	}

	// Generate destructor code for RHS symbols which are not used in the
	// reduce code
	for i := 0; i < rp.nrhs; i++ {
		if len(rp.rhsalias[i]) > 0 && !used[i] {
			ErrorMsg(lemp.filename, rp.ruleline,
				"Label %s for \"%s(%s)\" is never used.",
				rp.rhsalias[i], rp.rhs[i].name, rp.rhsalias[i])
			lemp.errorcnt++
		} else if len(rp.rhsalias[i]) == 0 {
			if has_destructor(rp.rhs[i], lemp) {
				fmt.Fprintf(out, "  yy_destructor(%d,&yymsp[%d].minor);\n",
					rp.rhs[i].index, i-rp.nrhs+1)
				(*lineno)++
			} else {
				fmt.Fprintf(out, "        /* No destructor defined for %s */\n",
					rp.rhs[i].name)
				(*lineno)++
			}
		}
	}

}

func has_destructor(sp *symbol, lemp *lemon) bool {
	if sp.typ == TERMINAL {
		return len(lemp.tokendest) > 0
	} else {
		return len(lemp.vardest) > 0 || len(sp.destructor) > 0
	}
}

/*
** Print the definition of the union used for the parser's data stack.
** This union contains fields for every possible data type for tokens
** and nonterminals.  In the process of computing and printing this
** union, also set the ".dtnum" field of every terminal and nonterminal
** symbol.
** 向语法分析器文件打印一种数据结构YYMINORTYPE，用来统一表示终结符和非终结符
*/
func print_stack_union(out *os.File, lemp *lemon, plineno *int, mhflag int) {

	lineno := *plineno

	// Allocate and initialize types[] and allocate stddt[]
	arraysize := lemp.nsymbol * 2
	types := make([]string, arraysize, arraysize)
	var stddt string         //  Standardized name for a datatype 数据类型的名称

	// Build a hash table of datatypes. The ".dtnum" field of each symbol
	// is filled in with the hash index plus 1.  A ".dtnum" value of 0 is
	// used for terminal symbols.  If there is no %default_type defined then
	// 0 is also used as the .dtnum value for nonterminals which do not specify
	// a datatype using the %type directive.
	for _, sp := range lemp.symbols {
		if sp == lemp.errsym {
			sp.dtnum = arraysize + 1
			continue
		}

		// 终结符或者没有指定全局数据类型且本身也没有特殊声明数据类型的非终结符 dtnum 为 0
		if sp.typ != NONTERMINAL || (len(sp.datatype) == 0 && len(lemp.vartype) == 0) {
			sp.dtnum = 0
			continue
		}

		cp := sp.datatype
		if len(cp) == 0 {
			cp = lemp.vartype
		}

		// 建立符号数据类型的哈希表
		cp = strings.TrimSpace(cp)
		stddt = cp
		hash := stringhash(stddt)
		hash = (hash & 0x7fffffff) % arraysize

		for len(types[hash]) > 0 { // 处理哈希碰撞
			if stddt == types[hash] {
				sp.dtnum = hash + 1
				break
			}
			hash ++
			if hash >= arraysize {
				hash = 0
			}
		}

		if len(types[hash]) == 0 {
			sp.dtnum = hash + 1
			types[hash] = stddt
		}
	}
	// 经过上面的处理对于每个非终结符都可以把dtnum-1作为hash值，然后从哈希表中获取到数据类型

	//  Print out the definition of YYTOKENTYPE and YYMINORTYPE
	name := "Parse"
	if len(lemp.name) > 0 {
		name = lemp.name
	}
	lineno = *plineno


	tokenType := "interface{}"
	if len(lemp.tokentype) > 0 {
		tokenType = lemp.tokentype
	}
	// TOKENTYPE为%token_type 指定的终结符的数据类型
	fmt.Fprintf(out, "type  %sTOKENTYPE %s\n", name, tokenType)
	lineno++

	fmt.Fprintf(out, "type  YYMINORTYPE struct {\n")
	lineno++
	fmt.Fprintf(out, "  yy0 %sTOKENTYPE\n", name)
	lineno++
	for i := 0; i < arraysize; i++ {
		if len(types[i]) == 0 {
			continue
		}
		fmt.Fprintf(out, "  yy%d %s\n", types[i], i+1)
		lineno++
	}

	fmt.Fprintf(out, "  yy%d int\n", lemp.errsym.dtnum)
	lineno++
	fmt.Fprintf(out, "} \n")
	lineno++
	*plineno = lineno
}

func stringhash(str string) int {
	hash := 0
	for _, char := range str {
		hash = hash*53 + int(char)
	}
	return hash
}

// Return the name of a datatype able to represent values between lwr and upr, inclusive.
// 用一个整数类型来表示所有的终结符
func minimum_size_type(lwr int, upr int) string {
	/*
	if lwr >= 0 {
		if upr < 255 {
			return "uint8"
		} else if upr > 65535 {
			return "uin16"
		} else {
			return "uint32"
		}
	} else if lwr >= -127 && upr <= 127 {
		return "int8"
	} else if lwr >= -32767 && upr < 32767 {
		return "int16"
	} else {
		return "int32"
	}
	*/
	return "int"
}

/*
** Each state contains a set of token transaction and a set of
** nonterminal transactions. Each of these sets makes an instance
** of the following structure.  An array of these structures is used
** to order the creation of entries in the yy_action[] table.
*/
type axset struct {
	stp     *state // A pointer to a state
	isTkn   bool   // True to use tokens.  False for non-terminals
	nAction int    /* Number of actions */
}

// Compare to axset structures for sorting purposes
func axset_compare(a *axset, b *axset) int {
	return b.nAction - a.nAction // 根据nAction大小逆序
}

type SortedAxset []*axset

func (sa SortedAxset) Len() int      { return len(sa) }
func (sa SortedAxset) Swap(i, j int) { sa[i], sa[j] = sa[j], sa[i] }
func (sa SortedAxset) Less(i, j int) bool {
	return axset_compare(sa[i], sa[j]) < 0
}

// Generate C source code for the parser
func ReportTable(lemp *lemon, mhflag int) { // mhflag: Output in makeheaders format if true

	in := tplt_open(lemp)
	if in == nil {
		return
	}

	out := file_open(lemp, ".go", os.O_WRONLY|os.O_CREATE)
	if out == nil {
		in.Close()
		return
	}

	lineno := 1
	tplt_xfer(lemp.name, in, out, &lineno) // TODO ???

	// Generate the include code, if any  打印语法文件中%include中的内容
	tplt_print(out, lemp, lemp.include, lemp.includeln, &lineno)

	tplt_xfer(lemp.name, in, out, &lineno)

	//  Generate #defines for all tokens
	// 打印所有终结符 类似 const PLUS   =  1
	if mhflag > 0 {
		prefix := ""
		if len(lemp.tokenprefix) > 0 {
			prefix = lemp.tokenprefix
		}

		for i := 1; i < lemp.nterminal; i++ {
			fmt.Fprintf(out, "const %s%-30s = %2d\n", prefix, lemp.symbols[i].name, i);
			lineno++;
		}
	}
	tplt_xfer(lemp.name, in, out, &lineno)

	// YYCODETYPE 用一个正数来表示终结符或非终结符
	fmt.Fprintf(out, "type YYCODETYPE %s\n",
		minimum_size_type(0, lemp.nsymbol+5));
	lineno++
	// YYNOCODE表示一个非法的符号
	fmt.Fprintf(out, "const YYNOCODE = %d\n", lemp.nsymbol+1);
	lineno++
	fmt.Fprintf(out, "type YYACTIONTYPE %s\n",
		minimum_size_type(0, lemp.nstate+lemp.nrule+5));
	lineno++
	print_stack_union(out, lemp, &lineno, mhflag)

	// 分析时使用堆栈大小，默认为100
	if len(lemp.stacksize) > 0 {
		if stacksize, err := strconv.Atoi(lemp.stacksize); err != nil || stacksize < 0 {
			ErrorMsg(lemp.filename, 0,
				"Illegal stack size: [%s].  The stack size should be an integer constant.",
				lemp.stacksize);
			lemp.errorcnt++;
			lemp.stacksize = "100"
		}
		fmt.Fprintf(out, "const YYSTACKDEPTH = %s\n", lemp.stacksize)
		lineno++
	} else {
		fmt.Fprintf(out, "const YYSTACKDEPTH = 100\n")
	}

	//lemp.arg = "Parse *pParse"
	var ARG_SDECL = "" // lemp.arg声明为变量
	var ARG_PDECL = "" // lemp.arg声明为函数中的参数
	var ARG_FETCH = "" // 把yypParser的pParse属性取出
	var ARG_STORE = "" // 设置yypParser的pParse属性
	if len(lemp.arg) > 0 { // TODO ???
		i := len(lemp.arg)
		for i >= 1 && util.IsSpace(lemp.arg[i-1]) {
			i --
		}

		for i >= 1 && (util.IsAlumn(lemp.arg[i-1]) || lemp.arg[i-1] == '_') {
			i --
		}
		ARG_SDECL = lemp.arg
		ARG_PDECL = fmt.Sprintf(",%s", lemp.arg)
		ARG_FETCH = fmt.Sprintf("%s = yypParser.%s", lemp.arg, lemp.arg[i:])     // TODO ???
		ARG_STORE = fmt.Sprintf("yypParser.%s = %s", lemp.arg[i:], lemp.arg[i:]) // TODO ???
	}

	fmt.Printf("ARG_SDECL: %s\n", ARG_SDECL)
	fmt.Printf("ARG_PDECL: %s\n", ARG_PDECL)
	fmt.Printf("ARG_FETCH: %s\n", ARG_FETCH)
	fmt.Printf("ARG_STORE: %s\n", ARG_STORE)


	fmt.Fprintf(out, "const YYNSTATE = %d\n", lemp.nstate)
	lineno++
	fmt.Fprintf(out, "const YYNRULE = %d\n", lemp.nrule)
	lineno++;
	fmt.Fprintf(out, "const YYERRORSYMBOL = %d\n", lemp.errsym.index)
	lineno++
	fmt.Fprintf(out, "const YYERRSYMDT = yy%d\n", lemp.errsym.dtnum)
	lineno++

	in.Close()
	out.Close()
}

const NO_OFFSET = -2147483647
