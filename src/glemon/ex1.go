package glemon

import (
	"unsafe"
	"fmt"
	"os"
)



//#line 12 "C:\code\lemon\ex1.go"
/* Next is all token values, in a form suitable for use by makeheaders.
** This section will be null unless lemon is run with the -m switch.
*/
/*
** These constants (all generated automatically by the parser generator)
** specify the various kinds of tokens (terminals) that the parser
** understands.
**
** Each symbol here is a terminal symbol in the grammar.
** 定义所有的终结符
*/
const PLUS                           =  1
const MINUS                          =  2
const DIVIDE                         =  3
const TIMES                          =  4
const INTEGER                        =  5
/* The next thing included is series of defines which control
** various aspects of the generated parser.
**    int         is the data type used for storing terminal
**                       and nonterminal numbers.  "unsigned char" is
**                       used if there are fewer than 250 terminals
**                       and nonterminals.  "int" is used otherwise.
**    YYNOCODE           is a number of type int which corresponds
**                       to no legal terminal or nonterminal number.  This
**                       number is used to fill in empty slots of the hash
**                       table.
**    YYFALLBACK         If defined, this indicates that one or more tokens
**                       have fall-back values which should be used if the
**                       original value of the token will not parse.
**    int       is the data type used for "action codes" - numbers
**                       that indicate what to do in response to the next
**                       token.
**    int     is the data type used for minor tokens given
**                       directly to the parser from the tokenizer.
**    YYMINORTYPE        is the data type used for all minor tokens.
**                       This is typically a union of many types, one of
**                       which is int.  The entry in the union
**                       for base tokens is called "yy0".
**    YYSTACKDEPTH       is the maximum depth of the parser's stack.
**         A static variable declaration for the %extra_argument
**         A parameter declaration for the %extra_argument
**         Code to store %extra_argument into yypParser
**         Code to extract %extra_argument from yypParser
**    YYNSTATE           the combined number of states.
**    YYNRULE            the number of rules in the grammar
**    YYERRORSYMBOL      is the code number of the error symbol.  If not
**                       defined, then do no error processing.
*/
 
const YYNOCODE = 10

type  YYMINORTYPE struct {
	yy0 int
	yy19 int
}
const YYSTACKDEPTH = 100
const YYNSTATE = 11
const YYNRULE = 6
const YYERRORSYMBOL = 6
const YY_NO_ACTION        =     YYNSTATE+YYNRULE+2
const YY_ACCEPT_ACTION    =     YYNSTATE+YYNRULE+1
const YY_ERROR_ACTION     =     YYNSTATE+YYNRULE

/* Next are that tables used to determine what action to take based on the
** current state and lookahead token.  These tables are used to implement
** functions that take a state number and lookahead value and return an
** action integer.
** 建立一个一维的yy_action[]数组，当输入表示状态state的整数和表示先行符号的整数，返回
** 表示相应动作的整数N。
**
** Suppose the action integer is N.  Then the action is determined as
** follows
**
**   0 <= N < YYNSTATE                  Shift N.  That is, push the lookahead
**                                      token onto the stack and goto state N.
**                                       表示移进动作
**
**   YYNSTATE <= N < YYNSTATE+YYNRULE   Reduce by rule N-YYNSTATE.
**                                       表示序号为N-YYNSTATE的产生式进行归约
**
**   N == YYNSTATE+YYNRULE              A syntax error has occurred.  语法错误
**
**   N == YYNSTATE+YYNRULE+1            The parser accepts its input. 接受状态
**
**   N == YYNSTATE+YYNRULE+2            No such action.  Denotes unused  不可能有的动作
**                                      slots in the yy_action[] table.
**
** The action table is constructed as a single large table named yy_action[].
** Given state S and lookahead X, the action is computed as
**
**      yy_action[ yy_shift_ofst[S] + X ]
**
** If the index value yy_shift_ofst[S]+X is out of range or if the value
** yy_lookahead[yy_shift_ofst[S]+X] is not equal to X or if yy_shift_ofst[S]
** is equal to YY_SHIFT_USE_DFLT, it means that the action is not in the table
** and that yy_default[S] should be used instead.
**
** The formula above is for computing the action when the lookahead is
** a terminal symbol.  If the lookahead is a non-terminal (as occurs after
** a reduce action) then the yy_reduce_ofst[] array is used in place of
** the yy_shift_ofst[] array and YY_REDUCE_USE_DFLT is used in place of
** YY_SHIFT_USE_DFLT.
**
** The following are the tables generated in this section:
**
**  yy_action[]        A single table containing all actions.
**  yy_lookahead[]     A table containing the lookahead for each entry in
**                     yy_action.  Used to detect hash collisions.
**  yy_shift_ofst[]    For each state, the offset into yy_action for
**                     shifting terminals.
**  yy_reduce_ofst[]   For each state, the offset into yy_action for
**                     shifting non-terminals after a reduce.
**  yy_default[]       Default action for each state.
*/
var yy_action  = []int {
	/*     0 */    11,    4,    2,    8,    6,   18,    1,    8,    6,    3,
	/*    10 */    10,    5,   17,   17,    7,    9,
}
var yy_lookahead =  []int {
	/*     0 */     0,    1,    2,    3,    4,    7,    8,    3,    4,    8,
	/*    10 */     5,    8,    9,    9,    8,    8,
}
const YY_SHIFT_USE_DFLT = -1
var  yy_shift_ofst  = []int {
	/*     0 */     5,    0,    5,    4,    5,    4,    5,   -1,    5,   -1,
	/*    10 */    -1,
}
const YY_REDUCE_USE_DFLT = -3
var yy_reduce_ofst  = []int {
	/*     0 */    -2,   -3,    1,   -3,    3,   -3,    6,   -3,    7,   -3,
	/*    10 */    -3,
}
var yy_default  = []int {
	/*     0 */    17,   17,   17,   12,   17,   13,   17,   14,   17,   15,
	/*    10 */    16,
}
var  YY_SZ_ACTTAB = len(yy_action)

/* The next table maps tokens into fallback tokens.  If a construct
** like the following:
**
**      %fallback ID X Y Z.
**
** appears in the grammer, then ID becomes a fallback token for X, Y,
** and Z.  Whenever one of the tokens X, Y, or Z is input to the parser
** but it does not parse, the type of the token is changed to ID and
** the parse is retried before an error is thrown.
*/

//#ifdef YYFALLBACK
var yyFallback = []int {
};
//#endif


/* The following structure represents a single element of the
** parser's stack.  Information stored includes:
**
**   +  The state number for the parser at this level of the stack.
**
**   +  The value of the token stored at this level of the stack.
**      (In other words, the "major" token.)
**
**   +  The semantic value stored at this level of the stack.  This is
**      the information used by the action routines in the grammar.
**      It is sometimes called the "minor" token.
*/
type yyStackEntry struct {
	stateno int       // The state-number
	major int         // The major token value.  This is the code
	// number for the token at this stack level 语法符号的序号
	minor YYMINORTYPE // The user-supplied minor token value.  This
	// is the value of the token 语法符号的值
}

/* The state of the parser is completely contained in an instance of
** the following structure */
type yyParser struct {
	yyidx int                     // Index of top element in stack 栈定的位置
	yyerrcnt int                  // Shifts left before out of the error
	// A place to hold %extra_argument
	yystack [YYSTACKDEPTH] yyStackEntry  // The parser's stack
}


var yyTokenName = []string {
	"$",             "PLUS",          "MINUS",         "DIVIDE",
	"TIMES",         "INTEGER",       "error",         "program",
	"expr",
}

var yyRuleName = []string{
	/*   0 */ "program ::= expr",
	/*   1 */ "expr ::= expr MINUS expr",
	/*   2 */ "expr ::= expr PLUS expr",
	/*   3 */ "expr ::= expr TIMES expr",
	/*   4 */ "expr ::= expr DIVIDE expr",
	/*   5 */ "expr ::= INTEGER",
}

/*
** This function returns the symbolic name associated with a token
** value.
*/
func ParseTokenName(tokenType int) string {
	return "";
}

/*
** This function allocates a new parser.
** The only argument is a pointer to a function which works like
** malloc.
**
** Inputs:
** A pointer to the function used to allocate memory.
**
** Outputs:
** A pointer to a parser.  This pointer is used in subsequent calls
** to Parse and ParseFree.
*/
func ParseAlloc( ) *yyParser{
	pParser := &yyParser{
		yyidx:-1,
	}
	return pParser
}


/*
** Pop the parser's stack once.
**
** If there is a destructor routine associated with the token which
** is popped from the stack, then call it.
**
** Return the major token number for the symbol popped.
*/
func  yy_pop_parser_stack(pParser *yyParser) int  {
	var  yymajor int //  var  yymajor int
	yytos := &pParser.yystack[pParser.yyidx]
	if pParser.yyidx<0 {
		return 0
	}
	yymajor = yytos.major
	//yy_destructor( yymajor, &yytos.minor)
	pParser.yyidx--
	return yymajor
}

/*
** Deallocate and destroy a parser.  Destructors are all called for
** all stack elements before shutting the parser down.
** 析构整个语法分析栈
** Inputs:
** <ul>
** <li>  A pointer to the parser.  This should be a pointer
**       obtained from ParseAlloc.
** <li>  A pointer to a function used to reclaim memory obtained
**       from malloc.
** </ul>
*/
func ParseFree(pointer unsafe.Pointer){
	pParser := (*yyParser) (pointer)
	if  pParser==nil  {
		return
	}
	for  pParser.yyidx >= 0 {
		yy_pop_parser_stack(pParser)
	}

}

/*
** Find the appropriate action for a parser given the terminal
** look-ahead token iLookAhead.
**
** If the look-ahead token is YYNOCODE, then check to see if the action is
** independent of the look-ahead.  If it is, return the action, otherwise
** return YY_NO_ACTION.
*/
func yy_find_shift_action(pParser *yyParser, iLookAhead int) int{
	var  i int
	stateno := pParser.yystack[pParser.yyidx].stateno
	i = yy_shift_ofst[stateno]
	if  i == YY_SHIFT_USE_DFLT  {
		return yy_default[stateno]
	}
	if iLookAhead == YYNOCODE  {
		return YY_NO_ACTION
	}
	i += iLookAhead
	if i<0 || i>=YY_SZ_ACTTAB || yy_lookahead[i] != int(iLookAhead) {
		var  iFallback   = yyFallback[iLookAhead]    // Fallback token
		if  iLookAhead < len(yyFallback)  && iFallback !=0  {
			return yy_find_shift_action(pParser, iFallback)
		}
		return yy_default[stateno]
	}else{
		return yy_action[i]
	}
}

/*
** Find the appropriate action for a parser given the non-terminal
** look-ahead token iLookAhead.
**
** If the look-ahead token is YYNOCODE, then check to see if the action is
** independent of the look-ahead.  If it is, return the action, otherwise
** return YY_NO_ACTION.
*/
func  yy_find_reduce_action(pParser *yyParser, iLookAhead int) int {
	var i int
	stateno := pParser.yystack[pParser.yyidx].stateno
	i = yy_reduce_ofst[stateno]
	if  i== YY_REDUCE_USE_DFLT  {
		return yy_default[stateno]
	}
	if  iLookAhead==YYNOCODE  {
		return YY_NO_ACTION
	}
	i += iLookAhead;
	if  i<0 || i>=YY_SZ_ACTTAB || yy_lookahead[i]!=iLookAhead  {
		return yy_default[stateno]
	}else{
		return yy_action[i]
	}
}

/*
** Perform a shift action.
*/
func yy_shift(yypParser *yyParser, yyNewState int, yyMajor int, yypMinor *YYMINORTYPE){
	var yytos *yyStackEntry
	yypParser.yyidx++;

	if yypParser.yyidx >= YYSTACKDEPTH  {

		yypParser.yyidx--

		for yypParser.yyidx >= 0   {
			yy_pop_parser_stack(yypParser)
		}
		// Here code is inserted which will execute if the parser
		// stack every overflows
		// Suppress warning about unused %extra_argument var
		return
	}
	yytos = &yypParser.yystack[yypParser.yyidx] // TODO nil 判断
	yytos.stateno = yyNewState
	yytos.major = yyMajor
	yytos.minor = *yypMinor
}

/* The following table contains information about every rule that
** is used during the reduce.
*/
type ruleInfo struct {
	lhs   int      /* Symbol on the left-hand side of the rule */
	nrhs  int   /* Number of right-hand side symbols in the rule */
}
/*
ruleInfo定义了产生式左边符号和右边符号的数目，在一个产生式右边所有符号都入栈后进行归约，
又要出栈，ruleInfo记录了要弹出多少个符号
*/
var yyRuleInfo  = []ruleInfo {
	{ 7, 1 },
	{ 8, 3 },
	{ 8, 3 },
	{ 8, 3 },
	{ 8, 3 },
	{ 8, 1 },
}


/*
** Perform a reduce action and the shift that must immediately
** follow the reduce.
*/
func yy_reduce(yypParser *yyParser, yyruleno int){
	var  yygoto int                     /* The next state */
	var  yyact int                      /* The next action */
	var  yygotominor YYMINORTYPE        /* The LHS of the rule reduced */
	var  yymsp *yyStackEntry            /* The top of the parser's stack */
	var  yysize int                     /* Amount to pop the stack */

	yymsp = &yypParser.yystack[yypParser.yyidx]

	switch  yyruleno  {
	/* Beginning here are the reduction cases.  A typical example follows:
	**   case 0:
	**  #line <lineno> <grammarfile>
	**     { ... }           // User supplied code
	**  #line <lineno> <thisfile>
	**     break;
	*/
	case 0:
		//#line 17 "C:\code\lemon\ex1.y"
		fmt.Printf("Result = %d\n", yypParser.yystack[yypParser.yyidx-0].minor.yy0)
	case 1:
		//#line 19 "C:\code\lemon\ex1.y"
		yygotominor.yy0 = yypParser.yystack[yypParser.yyidx-2].minor.yy0 - yypParser.yystack[yypParser.yyidx-0].minor.yy0;
	case 2:
		//#line 20 "C:\code\lemon\ex1.y"
		yygotominor.yy0 = yymsp[-2].minor.yy0 + yymsp[0].minor.yy0;
	case 3:
		//#line 21 "C:\code\lemon\ex1.y"
		yygotominor.yy0 = yymsp[-2].minor.yy0 * yymsp[0].minor.yy0;
	case 4:
		//#line 22 "C:\code\lemon\ex1.y"

		if(yymsp[0].minor.yy0 != 0) {
			yygotominor.yy0 = yymsp[-2].minor.yy0 / yymsp[0].minor.yy0;
		} else {
			fmt.Printf("Divide by zero!\n");
		}

	case 5:
		//#line 30 "C:\code\lemon\ex1.y"
		yygotominor.yy0 = yymsp[0].minor.yy0;
	}
	yygoto = yyRuleInfo[yyruleno].lhs
	yysize = yyRuleInfo[yyruleno].nrhs
	yypParser.yyidx -= yysize
	yyact = yy_find_reduce_action(yypParser,yygoto)
	if  yyact < YYNSTATE  {
		yy_shift(yypParser,yyact,yygoto,&yygotominor)
	}else if( yyact == YYNSTATE + YYNRULE + 1 ){
		yy_accept(yypParser)
	}
}

/*
** The following code executes when the parse fails
*/
func yy_parse_failed( yypParser *yyParser) {

	for yypParser.yyidx >= 0 {
		yy_pop_parser_stack(yypParser)
	}
	/* Here code is inserted which will be executed whenever the parser fails */
	/* Suppress warning about unused %extra_argument variable */
}

/*
** The following code executes when a syntax error first occurs.
*/

func yy_syntax_error( yypParser *yyParser, yymajor int, yyminor YYMINORTYPE) {

	//type TOKEN yyminor.yy0 // TODO
	//#line 12 "C:\code\lemon\ex1.y"

	fmt.Printf("Syntax error!\n");
	os.Exit(1)

	//#line 486 "C:\code\lemon\ex1.go"
	/* Suppress warning about unused %extra_argument variable */
}

/*
** The following is executed when the parser accepts
*/
func yy_accept(yypParser *yyParser){

	for yypParser.yyidx>=0 {
		yy_pop_parser_stack(yypParser)
	}
	/* Here code is inserted which will be executed whenever the parser accepts */
	; /* Suppress warning about unused %extra_argument variable */
}

/* The main parser program.
** The first argument is a pointer to a structure obtained from
** "ParseAlloc" which describes the current state of the parser.
** The second argument is the major token number.  The third is
** the minor token.  The fourth optional argument is whatever the
** user wants (and specified in the grammar) and is available for
** use by the action routines.
**
** Inputs:
** <ul>
** <li> A pointer to the parser (an opaque structure.)
** <li> The major token number.
** <li> The minor token number.
** <li> An option argument of a grammar-specified type.
** </ul>
**
** Outputs:
** None.
*/
func Parse(yyp unsafe.Pointer, yymajor int, yyminor int  ){
	var yyminorunion YYMINORTYPE
	var yyact int             /* The parser action. */
	var yyendofinput bool     /* True if we are at the end of input */
	var yyerrorhit   = false  /* True if yymajor has invoked an error */
	var yypParser = (*yyParser)(yyp) /* The parser */


	if yypParser.yyidx < 0  {
		if yymajor == 0 {
			return;
		}
		yypParser.yyidx = 0
		yypParser.yyerrcnt = -1
		yypParser.yystack[0].stateno = 0
		yypParser.yystack[0].major = 0
	}
	yyminorunion.yy0 = yyminor
	yyendofinput = (yymajor==0)


	for {
		yyact = yy_find_shift_action(yypParser, yymajor)
		if yyact < YYNSTATE  {
			yy_shift(yypParser, yyact, yymajor, &yyminorunion)
			yypParser.yyerrcnt--
			if  yyendofinput && yypParser.yyidx>=0 {
				yymajor = 0
			} else {
				yymajor = YYNOCODE
			}
		}else if yyact < YYNSTATE + YYNRULE  {
			yy_reduce(yypParser, yyact-YYNSTATE)
		}else if  yyact == YY_ERROR_ACTION  {
			var yymx int

			if YYERRORSYMBOL > 0 {
				/* A syntax error has occurred.
				** The response to an error depends upon whether or not the
				** grammar defines an error token "ERROR".
				**
				** This is what we do if the grammar does define ERROR:
				**
				**  * Call the %syntax_error function.
				**
				**  * Begin popping the stack until we enter a state where
				**    it is legal to shift the error symbol, then shift
				**    the error symbol.
				**
				**  * Set the error count to three.
				**
				**  * Begin accepting and shifting new tokens.  No new error
				**    processing will occur until three tokens have been
				**    shifted successfully.
				**
				*/
				if  yypParser.yyerrcnt < 0 {
					yy_syntax_error(yypParser, yymajor, yyminorunion)
				}
				yymx = yypParser.yystack[yypParser.yyidx].major
				if  yymx==YYERRORSYMBOL || yyerrorhit  {
					//yy_destructor(yymajor,&yyminorunion)
					yymajor = YYNOCODE;
				} else {
					for  {
						yyact = yy_find_shift_action(yypParser,YYERRORSYMBOL)
						if !(yypParser.yyidx >= 0 && yymx != YYERRORSYMBOL && yyact >= YYNSTATE) {
							break
						}
						yy_pop_parser_stack(yypParser)
					}

					if  yypParser.yyidx < 0 || yymajor==0 {
						//yy_destructor(yymajor,&yyminorunion)
						yy_parse_failed(yypParser)
						yymajor = YYNOCODE
					} else if  yymx!=YYERRORSYMBOL {
						 var u2 YYMINORTYPE
						u2.yy19 = 0
						yy_shift(yypParser,yyact,YYERRORSYMBOL, &u2)
					}
				}
				yypParser.yyerrcnt = 3
				yyerrorhit = 1

			} else {

				/* YYERRORSYMBOL is not defined */
				/* This is what we do if the grammar does not define ERROR:
				**
				**  * Report an error message, and throw away the input token.
				**
				**  * If the input token is $, then fail the parse.
				**
				** As before, subsequent error messages are suppressed until
				** three input tokens have been successfully shifted.
				*/
				if yypParser.yyerrcnt<=0 {
					yy_syntax_error(yypParser,yymajor,yyminorunion)
				}
				yypParser.yyerrcnt = 3
				//yy_destructor(yymajor,&yyminorunion)
				if yyendofinput {
					yy_parse_failed(yypParser)
				}
				yymajor = YYNOCODE
			}
		}else{
			yy_accept(yypParser)
			yymajor = YYNOCODE
		}

		if !(yymajor != YYNOCODE && yypParser.yyidx >= 0) {
			break
		}
	}
}
//#line 32 "C:\code\lemon\ex1.y"
 