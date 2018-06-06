package glemon

import (
	"os"
	"strings"
	"fmt"
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
	for i := 0; i < rp.nrhs; i++ {
		if i == cfp.dot {
			fmt.Fprintf(fp, " *")
		}
		/*if i == rp.nrhs {
			break;
		}*/
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
		fmt.Fprintf(fp, "%*s shift  %d", indent, ap.sp.name, ap.stp.index)
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

	out := file_open(lemp, "_def.go", os.O_WRONLY|os.O_APPEND)
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
//
func CompressTables(lemp *lemon) {
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

func ReportTable(lemp *lemon, mhflag int) {

}
