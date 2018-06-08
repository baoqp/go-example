package glemon

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

const TEST = false

// 命令参数
var (
	basisflag  *int
	compress   *int
	azDefine   *string
	rpflag     *int
	mhflag     *int
	quiet      *int
	statistics *int
	version    *int
)

func ISOPT(x string) bool {
	return x[0] == '-' || x[0] == '+' || x[0] == '='
}

func main() {
	basisflag = flag.Int("b", 0, "Pr*int only the basis in report.")
	compress = flag.Int("c", 0, "Don't compress the action table.")
	azDefine = flag.String("D", "", "Define an %ifdef macro.")
	rpflag = flag.Int("g", 0, "Print grammar without actions.")
	mhflag = flag.Int("m", 0, "Output a makeheaders compatible file.")
	quiet = flag.Int("q", 0, "Quiet) Don't print the report file.")
	statistics = flag.Int("s", 0, "Print parser stats to standard output.")
	version = flag.Int("x", 0, "Print the version number.")
	flag.Parse()

	if *version > 0 {
		fmt.Println("Gemon version 1.0")
		os.Exit(0)
	}

	args := os.Args[1:]
	filename := args[len(args)-1] // 要处理的文件
	filename = "C:\\personal\\lemon\\ex1.y"
	if ISOPT(filename) {
		fmt.Println("no file present")
		os.Exit(1)
	}

	Strsafe_init()
	Symbol_init()
	State_init()

	lem := lemon{}
	lem.filename = filename
	lem.basisflag = *basisflag

	Symbol_new("$")
	lem.errsym = Symbol_new("error")
	Parse(&lem)
	if lem.errorcnt > 0 {
		os.Exit(lem.errorcnt)
	}
	if lem.rule == nil {
		fmt.Println("Empty grammar.")
		os.Exit(1)
	}

	//  Count and index the symbols of the grammar
	lem.nsymbol = Symbol_count()
	Symbol_new("{default}")
	lem.symbols = Symbol_arrayof()
	for i := 0; i <= lem.nsymbol; i++ {
		lem.symbols[i].index = i
	}

	sort.Sort(SortedSymol(lem.symbols))

	lem.nterminal = 0
	for i := 0; i <= lem.nsymbol; i++ {
		lem.symbols[i].index = i
		if lem.symbols[i].typ == TERMINAL {
			lem.nterminal ++
		}

	}

	//  Initialize the size for all follow and first sets
	SetSize(lem.nterminal + 1)

	// Find the precedence for every production rule (that has one)
	FindRulePrecedences(&lem)

	// Compute the lambda-nonterminals and the first-sets for every nonterminal
	FindFirstSets(&lem)

	/*
	fmt.Println("---println first set---")
	for _, symbol := range lem.symbols {
		fmt.Printf("%s %d, %v \n", symbol.name, symbol.typ, symbol.firstset)
	}
	*/

	// Compute all LR(0) states.  Also record follow-set propagation
	// links so that the follow-set can be computed later
	lem.nstate = 0
	FindStates(&lem)
	lem.sorted = State_arrayof()

	// Tie up loose ends on the propagation links
	FindLinks(&lem)

	// Compute the follow set of every reducible configuration
	FindFollowSets(&lem)

	/*fmt.Println("---println fws---")
	for i := 0; i < lem.nstate; i++ {
		stp := lem.sorted[i]
		fmt.Printf("State %d:\n", stp.index)
		cfp := stp.cfp
		for cfp != nil {
			if cfp.dot == cfp.rp.nrhs {
				buf := fmt.Sprintf("(%d)", cfp.rp.index)
				fmt.Printf("    %5s ", buf)
			} else {
				fmt.Printf("          ")
			}
			PrintConfig(cfp)
			fmt.Printf("\n")
			cfp = cfp.next
		}
		fmt.Printf("\n")
	}
	*/

	// Compute the action tables
	FindActions(&lem);
	if *compress == 0 {
		CompressTables(&lem)
	}

	// Generate a report of the parser generated.  (the "y.output" file)
	if *quiet == 1 {
		ReportOutput(&lem)
	}

	//-------------------OK tag-------------------------------//
	// Generate the source code for the parser
	 ReportTable(&lem, 1)
}

// print config and fws
func PrintConfig(cfp *config) {
	rp := cfp.rp
	fmt.Printf("%s ::=", rp.lhs.name)
	for i := 0; i <= rp.nrhs; i++ {
		if i == cfp.dot {
			fmt.Printf(" *")
		}
		if i == rp.nrhs {
			break;
		}
		fmt.Printf(" %s", rp.rhs[i].name)
	}

	fmt.Printf(" %v", cfp.fws)
}
