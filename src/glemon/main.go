package glemon

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"util"
)

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
	//filename = "C://test.y"
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
	for i := 0; i <= lem.nsymbol; i++ {
		lem.symbols[i].index = i
	}

	fmt.Println("after sorted")
	for i := 1; util.IsUpperChar(lem.symbols[i].name[0]); i++ {
		lem.nterminal = i
	}

}
