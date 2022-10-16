package interval_test

import (
	"fmt"

	"github.com/rdleal/intervalst/interval"
)

func Example() {
	cmpFn := func(x, y int) int { return x - y }

	// We must define the value type as the compiler can't infer in this case.
	st := interval.NewSearchTree[string](cmpFn)

	st.Insert(17, 19, "value1")
	st.Insert(5, 8, "value2")
	st.Insert(21, 24, "value3")
	st.Insert(4, 8, "value4")
	st.Insert(15, 18, "value5")
	st.Insert(7, 10, "value6")
	st.Insert(16, 22, "value7")

	val, ok := st.AnyIntersection(23, 25)
	fmt.Println(val, ok)

	val, ok = st.AnyIntersection(21, 23)
	fmt.Println(val, ok)

	_, ok = st.AnyIntersection(12, 14)
	fmt.Println(ok)

	// Output:
	// value3 true
	// value7 true
	// false
}
