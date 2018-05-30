package glemon

import "util"

//--------------------------s_x1--------------------------------//

type s_x1 map[string]string

var x1a s_x1

func Strsafe_init() {
	tmp := make(map[string]string)
	x1a = s_x1(tmp)
}

func Strsafe_insert(data string) bool {
	if _, ok := x1a[data]; ok {
		/* An existing entry with the same key is found. */
		/* Fail because overwrite is not allows. */
		return false
	}

	x1a[data] = data
	return true
}

func Strsafe_find(key string) string {
	if _, ok := x1a[key]; ok {
		return x1a[key]
	}
	return ""
}

//----------------------------s_x2------------------------------//

/* Return a pointer to the (terminal or nonterminal) symbol "x".
** Create a new symbol if this is the first time "x" has been seen.
** 使用字符串x创建符号
*/
func Symbol_new(x string) *symbol {
	sp := Symbol_find(x)
	if sp == nil {
		sp = &symbol{
			name:       x,
			rule:       nil,
			fallback:   nil,
			prec:       -1,
			assoc:      UNK,
			firstset:   nil,
			lambda:     false,
			destructor: nil,
			datatype:   nil,
		}

		if util.IsUpper(x) {
			sp.typ = TERMINAL
		} else {
			sp.typ = NONTERMINAL
		}

		Symbol_insert(x, sp);
	}
	return sp;
}

/* Compare two symbols for working purposes
**
** Symbols that begin with upper case letters (terminals or tokens)
** must sort before symbols that begin with lower case letters
** (non-terminals).  Other than that, the order does not matter.
**
** We find experimentally that leaving the symbols in their original
** order (the order they appeared in the grammar file) gives the
** smallest parser tables in SQLite.*
** 先按是否为终结符排序，相同类型按照出现顺序排序
*/
func Symbolcmpp(a [][]*symbol, b [][]*symbol) int {

	c1 := 0
	if util.IsLower(a[0][0].name) {
		c1 = 1
	}

	c2 := 0
	if util.IsLower(b[0][0].name) {
		c2 = 1
	}

	i1 := a[0][0].index + 10000000*c1
	i2 := b[0][0].index + 10000000*c2

	return i1 - i2
}

type s_x2 map[string]*symbol

var x2a s_x2

func Symbol_init() {
	if x2a != nil {
		return
	}
	tmp := make(map[string]*symbol)
	x2a = s_x2(tmp)
}

func Symbol_insert(key string, data *symbol) bool {
	if _, ok := x2a[key]; ok {
		/* An existing entry with the same key is found. */
		/* Fail because overwrite is not allows. */
		return false
	}

	x2a[key] = data
	return true
}

func Symbol_find(key string) *symbol {
	if _, ok := x2a[key]; ok {
		return x2a[key]
	}
	return nil
}

func Symbol_count() int {
	return len(x2a)
}

func Symbol_arrayof() []*symbol {
	if x2a == nil {
		return nil
	}

	v := make([]*symbol, 0, len(x2a))
	for _, value := range x2a {
		v = append(v, value)
	}
	return v
}

//----------------------------s_x3------------------------------//

/* Compare two configurations */
func Configcmp(a *config, b *config) int {
	x := a.rp.index - b.rp.index
	if x == 0 {
		x = a.dot - b.dot
	}
	return x
}

/* Compare two states */
func statecmp(a *config, b *config) int {
	var rc int
	for rc = 0; rc == 0 && a != nil && b != nil; {
		rc = a.rp.index - b.rp.index
		a = a.bp
		b = b.bp
	}

	if rc == 0 {
		if a != nil {
			rc = 1
		}

		if b != nil {
			rc = -1
		}
	}

	return rc
}

func statehash(a *config) int {
	h := 0
	for a != nil {
		h = h*571 + a.rp.index*37 + a.dot;
		a = a.bp;
	}
	return h;
}

func State_new() *state {
	return &state{}
}

type s_x3node struct {
	data *state
	key  *config
}

// golang 的map不支持自定义key
type s_x3 struct {
	size  int           /* The number of available slots. Must be a power of 2 greater than or equal to 1 */
	count int           /* Number of currently slots filled */
	tbl   []*s_x3node   /* The data stored here */
	ht    [][]*s_x3node /* Hash table for lookups 保存指针到hash table加快查找*/
}

var x3a *s_x3

func State_init() {

	if x3a != nil {
		return
	}

	x3a = &s_x3{}
	x3a.size = 128
	x3a.count = 0
	x3a.tbl = make([]*s_x3node, 0, 128)
	x3a.ht = make([][]*s_x3node, 128, 128)

	// TODO 初始化
}

// To be test
func State_insert(data *state, key *config) bool {

	if x3a == nil {
		return false
	}
	ph := statehash(key)
	h := ph & (x3a.size - 1)
	np := x3a.ht[h]

	for _, node := range np {
		if node.key == key {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return false
		}
	}

	// 扩容
	if x3a.count == x3a.size {
		size := x3a.size * 2
		arr := make([]*s_x3node, x3a.size, size)
		copy(arr, x3a.tbl)
		x3a.ht = make([][]*s_x3node, size, size)
		for i := 0; i < x3a.count; i++ {
			h = statehash(x3a.tbl[i].key) & (size - 1)
			x3a.ht[h] = append(x3a.ht[h], x3a.tbl[i])
		}
		x3a.tbl = arr
		x3a.size = size
	}

	newstate := &s_x3node{
		key:  key,
		data: data,
	}
	x3a.tbl = append(x3a.tbl, newstate)
	h = ph & (x3a.size - 1)
	x3a.ht[h] = append(x3a.ht[h], newstate)
	x3a.count += 1
	return true
}

func State_find(key *config) *state {
	if x3a == nil {
		return nil
	}
	h := statehash(key) & (x3a.size - 1)
	np := x3a.ht[h]

	for _, node := range np {
		if statecmp(node.key, key) == 0 {
			return node.data
		}
	}
	return nil
}
