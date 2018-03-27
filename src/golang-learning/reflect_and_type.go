package golang_learning

import (
	"fmt"

	"os"
	"time"
	"math/rand"
	"strings"
	"unsafe"
	"reflect"
)


var types = []string{
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"byte", "rune", "uintptr", "bool", "string",
	"float32", "float64", "complex64", "complex128",
	"[]byte", "[]string", "map[string]int",
	"chan int", "func(int) int",
}

var values = []string{
	"0", "0", "0", "0", "0",
	"0", "0", "0", "0", "0",
	"0", "0", "0", "false", "\"\"",
	"0", "0", "0+0i", "0+0i",
	"[]byte{}", "[]string{}", "map[string]int{}",
	"nil", "func(int) int {return 0}",
}

const template1 = `package golang_learning

import (
	"fmt"
	"reflect"
)

var v = struct {
		a %v
		b %v
		c %v
		d %v
		e %v
}{%v, %v, %v, %v, %v}
`
const template2 = `
func init() {
	fmt.Printf("%#T\n", v)
    t := reflect.TypeOf(v)
    fmt.Printf("结构体大小：%v\n", t.Size())
    for i := 0; i < t.NumField(); i++ {
        showAlign(t, i)
    }
}

func showAlign(v reflect.Type, i int) {
    sf := v.Field(i)
    fmt.Printf("字段 %10v，大小：%2v，对齐：%2v，字段对齐：%2v，偏移：%2v\n",
        sf.Type.Kind(),
        sf.Type.Size(),
        sf.Type.Align(),
        sf.Type.FieldAlign(),
        sf.Offset,
    )
}`



func GoAlign() {
	f, err := os.OpenFile("testAlign.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	t := [5]string{}
	v := [5]string{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		n := rand.Intn(len(types))
		t[i] = types[n]
		v[i] = values[n]
	}

	fmt.Fprintf(f, template1,
		t[0], t[1], t[2], t[3], t[4],
		v[0], v[1], v[2], v[3], v[4],
	)
	fmt.Fprint(f, template2)
}

func AlterPrivateField() {
	// 创建一个 strings 包中的 Reader 对象
	// 它有三个私有字段：s string、i int64、prevRune int
	sr := strings.NewReader("abcdef")
	// 此时 sr 中的成员是无法修改的
	fmt.Println(sr) // &{abcdef 0 -1}
	// 但是我们可以通过 unsafe 来进行修改
	// 先将其转换为通用指针
	p := unsafe.Pointer(sr)
	// 获取结构体地址
	up0 := uintptr(p)
	// 确定要修改的字段（这里不能用 unsafe.Offsetof 获取偏移量，因为是私有字段）
	if sf, ok := reflect.TypeOf(*sr).FieldByName("i"); ok {
		// 偏移到指定字段的地址
		up := up0 + sf.Offset
		// 转换为通用指针
		p = unsafe.Pointer(up)
		// 转换为相应类型的指针
		pi := (*int64)(p)
		// 对指针所指向的内容进行修改
		*pi = 3 // 修改索引
	}
	// 看看修改结果
	fmt.Println(sr) //&{abcdef 3 -1}
	// 看看读出的是什么
	b, err := sr.ReadByte()
	fmt.Printf("%c, %v\n", b, err) // d, <nil>
}



// 获取一个对象的字段和方法
func GetMembers(i interface{}) {
	// 获取 i 的类型信息
	t := reflect.TypeOf(i)

	for {
		// 进一步获取 i 的类别信息
		if t.Kind() == reflect.Struct {
			// 只有结构体可以获取其字段信息
			fmt.Printf("\n%-8v %v 个字段:\n", t, t.NumField())
			// 进一步获取 i 的字段信息
			for i := 0; i < t.NumField(); i++ {
				fmt.Println(t.Field(i).Name)
			}
		}
		// 任何类型都可以获取其方法信息
		fmt.Printf("\n%-8v %v 个方法:\n", t, t.NumMethod())
		// 进一步获取 i 的方法信息
		for i := 0; i < t.NumMethod(); i++ {
			fmt.Println(t.Method(i).Name)
		}
		if t.Kind() == reflect.Ptr {
			// 如果是指针，则获取其所指向的元素
			t = t.Elem()
		} else {
			// 否则上面已经处理过了，直接退出循环
			break
		}
	}
}


// 定义一个结构体用来进行测试
type sr struct {
	string
}

// 接收器为实际类型
func (s sr) Read() {
}

// 接收器为指针类型
func (s *sr) Write() {
}