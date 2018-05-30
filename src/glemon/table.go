package glemon

// 计算字符串hash值
func strhash(x string) int {
	h := 0;
	for _, c := range x {
		h = h*13 + int(c)
	}
	return h;
}

/* There is one instance of the following structure for each
** associative array of type "x1".
** TODO hashmap ??
*/
type s_x1 struct {
	size  int          /* The number of available slots. Must be a power of 2 greater than or equal to 1 */
	count int          /* Number of currently slots filled */
	tbl   []s_x1node;  /* The data stored here */
	ht    []*s_x1node; /* Hash table for lookups */
}

/* There is only one instance of the array, which is the following */
var x1a *s_x1 = nil

/* Allocate a new associative array */
func Strsafe_init() {
	if x1a != nil {
		return
	}

	x1a = &s_x1{}
	x1a.size = 1024
	x1a.count = 0
	x1a.tbl = make([]s_x1node, 1024, 1024)
	x1a.ht = make([]*s_x1node, 1024, 1024)

	// TODO 初始化
}

/* There is one instance of this structure for every data element
** in an associative array of type "x1".
*/
type s_x1node struct {
	data string       /* The data */
	next *s_x1node;   /* Next entry with the same hash */
	from []*s_x1node; /* Previous link TODO ??? */
}

/* There is one instance of the following structure for each
** associative array of type "x2".
*/
type s_x2 struct {
	size  int            /* The number of available slots. Must be a power of 2 greater than or equal to 1 */
	count int            /* Number of currently slots filled */
	tbl   []s_x2node;    /* The data stored here */
	ht    [][]*s_x2node; /* Hash table for lookups */
}

/* There is only one instance of the array, which is the following */
var x2a *s_x2 = nil

/* Allocate a new associative array */
func Symbol_init() {
	if x2a != nil {
		return
	}

	x2a = &s_x2{}
	x2a.size = 128
	x2a.count = 0
	x2a.tbl = make([]s_x2node, 128, 128)
	x2a.ht = make([][]*s_x2node, 128, 128)

	// TODO 初始化
}

/* Return a pointer to data assigned to the given key.  Return NULL
** if no such key. */
func Symbol_find(key string) *symbol {
	if x2a == nil {
		return nil
	}
	h := strhash(key) & (x2a.size - 1)
	np := x2a.ht[h]

	for _, node := range np {
		if node.key == key {
			return node.data
		}
	}
	return nil
}

func Symbol_insert(data *symbol, key string) bool {

	if x2a == nil {
		return false
	}
	ph := strhash(key);
	h := ph & (x2a.size - 1)
	np := x2a.ht[h]

	for _, node := range np {
		if node.key == key {
			/* An existing entry with the same key is found. */
			/* Fail because overwrite is not allows. */
			return false
		}
	}

	// 扩容
	if x2a.count == x2a.size {

	}

	h = ph & (x2a.size-1);
	npp := &x2a.tbl[x2a.count]
	npp.key = key
	npp.data = data
	// TODO ???

}

/* There is one instance of this structure for every data element
** in an associative array of type "x2".
*/
type s_x2node struct {
	data *symbol      /* The data */
	key  string       /* The key */
	next *s_x1node;   /* Next entry with the same hash */
	from []*s_x1node; /* Previous link TODO ??? */
}

/* There is one instance of the following structure for each
** associative array of type "x2".
*/
type s_x3 struct {
	size  int          /* The number of available slots. Must be a power of 2 greater than or equal to 1 */
	count int          /* Number of currently slots filled */
	tbl   []s_x3node;  /* The data stored here */
	ht    []*s_x3node; /* Hash table for lookups */
}

/* There is only one instance of the array, which is the following */
var x3a *s_x3 = nil

/* Allocate a new associative array */
func State_init() {
	if x3a != nil {
		return
	}

	x3a = &s_x3{}
	x3a.size = 128
	x3a.count = 0
	x3a.tbl = make([]s_x3node, 128, 128)
	x3a.ht = make([]*s_x3node, 128, 123)

	// TODO 初始化
}

/* There is one instance of this structure for every data element
** in an associative array of type "x2".
*/
type s_x3node struct {
	data *state       /* The data */
	key  string       /* The key */
	next *s_x3node;   /* Next entry with the same hash */
	from []*s_x3node; /* Previous link TODO ??? */
}
