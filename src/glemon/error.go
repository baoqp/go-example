package glemon

import "fmt"

const (
	ERRMSGSIZE  = 10000
	LINEWIDTH   = 79
	PREFIXLIMIT = 30
)

func ErrorMsg(filename string, lineno int, format ...string) {
	var prefix string
	if lineno > 0 {
		prefix = fmt.Sprintf("%.*s:%d: ", PREFIXLIMIT-10, filename, lineno)
	} else {
		prefix = fmt.Sprintf("%.*s: ", PREFIXLIMIT-10, filename);
	}

	availablewidth := LINEWIDTH - len(prefix)

	var errmsg string
	for _, str := range format {
		errmsg += str
	}

	errmsgsize := len(errmsg)
	for ; errmsgsize > 0 && errmsg[errmsgsize-1] == '\n'; {
		errmsgsize -= 1
		errmsg = errmsg[:errmsgsize]
	}

	base := 0
	var end, restart int
	chars := []byte(errmsg)
	for ; base < len(chars); {
		restart = findbreak(chars[base:], 0, availablewidth)
		end = restart
		for ; chars[restart] == ' '; {
			restart ++
		}
		fmt.Printf("%s%.*s\n", prefix, end, string(chars[base:]))
	}

}

/* Find a good place to break "msg" so that its length is at least "min"
** but no more than "max".  Make the point as close to max as possible.
*/
func findbreak(msg []byte, min, max int) int {
	i := min
	spot := min
	var c byte
	for ; i <= max; i++ {
		c = msg[i]
		if c == '\t' {
			msg[i] = ' '
		}
		if c == '\n' {
			msg[i] = ' ';
			spot = i;
			break;
		}
		if c == 0 {
			spot = i;
			break
		}
		if c == '-' && i < max-1 {
			spot = i + 1
		}
		if c == ' ' {
			spot = i
		}
	}
	return spot
}
