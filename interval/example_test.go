package interval_test

import (
	"fmt"
	"time"

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

func ExampleSearchTree_Ceil() {
	cmpFn := func(x, y int) int { return x - y }

	st := interval.NewSearchTree[string](cmpFn)

	st.Insert(17, 19, "value1")
	st.Insert(5, 8, "value2")
	st.Insert(21, 24, "value3")
	st.Insert(4, 8, "value4")
	st.Insert(15, 18, "value5")
	st.Insert(7, 10, "value6")

	val, ok := st.Ceil(9, 16)
	fmt.Println(val, ok)
	// Output:
	// value5 true
}

func ExampleSearchTree_Floor() {
	cmpFn := func(x, y int) int { return x - y }

	st := interval.NewSearchTree[string](cmpFn)

	st.Insert(17, 19, "value1")
	st.Insert(5, 8, "value2")
	st.Insert(21, 24, "value3")
	st.Insert(4, 8, "value4")
	st.Insert(15, 18, "value5")
	st.Insert(7, 10, "value6")

	val, ok := st.Floor(9, 16)
	fmt.Println(val, ok)
	// Output:
	// value6 true
}

func ExampleMultiValueSearchTree_Insert() {
	cmpFn := func(start, end time.Time) int {
		switch {
		case start.After(end):
			return 1
		case start.Before(end):
			return -1
		default:
			return 0
		}
	}

	st := interval.NewMultiValueSearchTree[string](cmpFn)

	start, end := time.Now(), time.Now().Add(time.Hour)
	st.Insert(start, end, "event1", "event2", "event3")

	st.Insert(start, end, "event4")

	vals, ok := st.Find(start, end)
	fmt.Println(vals, ok)
	// Output:
	// [event1 event2 event3 event4] true
}

func ExampleMultiValueSearchTree_Upsert() {
	cmpFn := func(start, end time.Time) int {
		switch {
		case start.After(end):
			return 1
		case start.Before(end):
			return -1
		default:
			return 0
		}
	}

	st := interval.NewMultiValueSearchTree[string](cmpFn)

	start, end := time.Now(), time.Now().Add(time.Hour)
	st.Insert(start, end, "event1", "event2", "event3")

	// Upsert will replace the previous values associated
	// with the start and end interval key, if any.
	st.Upsert(start, end, "event4", "event5")

	vals, ok := st.Find(start, end)
	fmt.Println(vals, ok)
	// Output:
	// [event4 event5] true
}

func ExampleTreeWithIntervalPoint() {
	cmpFn := func(start, end time.Time) int {
		switch {
		case start.After(end):
			return 1
		case start.Before(end):
			return -1
		default:
			return 0
		}
	}

	st := interval.NewSearchTreeWithOptions[string](cmpFn, interval.TreeWithIntervalPoint())

	pointInerval := time.Now()
	st.Insert(pointInerval, pointInerval, "event")

	vals, ok := st.Find(pointInerval, pointInerval)
	fmt.Println(vals, ok)
	// Output:
	// event true
}
