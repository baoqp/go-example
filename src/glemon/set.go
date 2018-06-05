package glemon

/*
** Set manipulation routines for the LEMON parser generator.
*/

var size int = 0

func SetSize(n int) {
	size = n + 1
}

// size为终结符数量加1
func SetNew() []byte {
	set := make([]byte, size, size)
	for i := range set {
		set[i] = 0
	}
	return set
}

func SetFree(s []byte) {
	s = s[:0]
}

/* Add a new element to the set.  Return TRUE if the element was added
** and FALSE if it was already there. */
func SetAdd(s []byte, e int) int {
	rv := s[e]
	s[e] = 1
	if rv == 0 {
		return 1
	}
	return 0
}

// 把s2合并到s1中
func SetUnion(s1, s2 []byte) int {
	progress := 0
	for i := 0; i < size; i++ {
		if s2[i] == 0 {
			continue
		}

		if s1[i] == 0 {
			progress = 1
			s1[i] = 1
		}
	}

	return progress
}
