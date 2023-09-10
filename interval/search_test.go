package interval

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var timeCmp = func(start, end time.Time) int {
	switch {
	case start.After(end):
		return 1
	case start.Before(end):
		return -1
	default:
		return 0
	}
}

func TestSearchTree_AnyIntersection_Time(t *testing.T) {
	t.Run("HasIntersection", func(t *testing.T) {
		st := NewSearchTree[string](timeCmp)
		defer mustBeValidTree(t, st)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end, "date1")

		start, end = end.Add(1*time.Hour), end.Add(2*time.Hour)
		st.Insert(start, end, "date2")

		start, end = time.Now().Add(-(5 * time.Hour)), time.Now().Add(-(3 * time.Hour))
		st.Insert(start, end, "date3")

		start, end = start.Add(1*time.Hour), end.Add(1*time.Hour)

		got, ok := st.AnyIntersection(start, end)
		if !ok {
			t.Errorf("st.AnyIntersection(%v, %v): got no intersection", start, end)
		}

		if want := "date3"; got != want {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected value %v; want %v", start, end, got, want)
		}
	})

	t.Run("HasExactIntersection", func(t *testing.T) {
		st := NewSearchTree[int](timeCmp)
		defer mustBeValidTree(t, st)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end, 0)

		start, end = start.Add(2*time.Hour), end.Add(1*time.Hour)
		st.Insert(start, end, 1)

		start, end = start.Add(-(5 * time.Hour)), end.Add(-(3 * time.Hour))
		st.Insert(start, end, 2)

		got, ok := st.AnyIntersection(start, end)
		if !ok {
			t.Errorf("st.AnyIntersection(%v, %v): got no intersection", start, end)
		}

		if want := 2; got != want {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected value %v; want %v", start, end, got, want)
		}
	})

	t.Run("HasNoIntersection", func(t *testing.T) {
		st := NewSearchTree[float64](timeCmp)
		defer mustBeValidTree(t, st)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end, 0.0)

		start, end = start.Add(2*time.Hour), end.Add(1*time.Hour)
		st.Insert(start, end, 1.0)

		start, end = start.Add(5*time.Hour), end.Add(3*time.Hour)
		st.Insert(start, end, 2.0)

		start, end = start.Add(1*time.Hour), end.Add(1*time.Hour)

		got, ok := st.AnyIntersection(start, end)
		if ok {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected value: %v", start, end, got)
		}
	})
}

func TestSearchTree_AnyIntersection(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")
	st.Insert(16, 22, "node7")

	testCases := []struct {
		start   int
		end     int
		wantOK  bool
		wantVal string
	}{
		{
			start:   23,
			end:     25,
			wantOK:  true,
			wantVal: "node3",
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
		{
			start:   21,
			end:     23,
			wantOK:  true,
			wantVal: "node7",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.AnyIntersection(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.AnyIntersection(%v, %v): got intersection %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if got != tc.wantVal {
				t.Errorf("st.AnyIntersection(%v, %v): got unexpected interval %v; want %v", tc.start, tc.end, got, tc.wantVal)
			}
		})
	}
}

func TestSearchTree_AnyIntersection_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.AnyIntersection(1, 10)
	if ok {
		t.Errorf("st.AnyIntersect(1, 10): got unexpected value %v", got)
	}
}

func TestMultiValueSearchTree_AnyIntersection(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")
	st.Insert(16, 22, "node7")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    23,
			end:      25,
			wantOK:   true,
			wantVals: []string{"node3"},
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
		{
			start:    21,
			end:      23,
			wantOK:   true,
			wantVals: []string{"node7"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.AnyIntersection(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.AnyIntersection(%v, %v): got intersection %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.AnyIntersection(%v, %v): got unexpected interval value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestMultiValueSearchTree_AnyIntersection_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.AnyIntersection(1, 10)
	if ok {
		t.Errorf("st.AnyIntersect(1, 10): got unexpected value %v", got)
	}
}

func TestSearchTree_Find(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")
	st.Insert(16, 22, "node7")

	testCases := []struct {
		start   int
		end     int
		wantOK  bool
		wantVal string
	}{
		{
			start:   4,
			end:     8,
			wantOK:  true,
			wantVal: "node4",
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
		{
			start:   7,
			end:     10,
			wantOK:  true,
			wantVal: "node6",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Find(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Find(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if got != tc.wantVal {
				t.Errorf("st.Find(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVal)
			}
		})
	}
}

func TestMultiValueSearchTree_Find(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4", "node4.1")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")
	st.Insert(16, 22, "node7")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    4,
			end:      8,
			wantOK:   true,
			wantVals: []string{"node4", "node4.1"},
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
		{
			start:    7,
			end:      10,
			wantOK:   true,
			wantVals: []string{"node6"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Find(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Find(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.Find(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestSearchTree_AllIntersections(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    9,
			end:      16,
			wantOK:   true,
			wantVals: []string{"node6", "node5"},
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.AllIntersections(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.AllIntersections(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.AllIntersections(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestSearchTree_AllIntersections_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.AllIntersections(1, 10)
	if ok {
		t.Errorf("st.AllIntersections(1, 10): got unexpected value %v", got)
	}
}

func TestMultiValueSearchTree_AllIntersections(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6", "node6.1")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    9,
			end:      16,
			wantOK:   true,
			wantVals: []string{"node6", "node6.1", "node5"},
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.AllIntersections(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.AllIntersections(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.AllIntersections(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestMultiValueSearchTree_AllIntersections_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.AllIntersections(1, 10)
	if ok {
		t.Errorf("st.AllIntersections(1, 10): got unexpected value %v", got)
	}
}

func TestSearchTree_Min(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	want := "node4"

	got, ok := st.Min()
	if !ok {
		t.Error("st.Min(): got no min value")
	}

	if got != want {
		t.Errorf("st.Min(): got unexpected value %v; want %v", got, want)
	}
}

func TestSearchTree_Min_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.Min()
	if ok {
		t.Errorf("st.Min(): got unexpected min value %v", got)
	}
}

func TestMultiValueSearchTree_Min(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4", "node4.1")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	want := []string{"node4", "node4.1"}

	got, ok := st.Min()
	if !ok {
		t.Error("st.Min(): got no min value")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("st.Min(): got unexpected value %v; want %v", got, want)
	}
}

func TestMultiValueSearchTree_Min_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.Min()
	if ok {
		t.Errorf("st.Min(): got unexpected min value %v", got)
	}
}

func TestSearchTree_Max(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	want := "node3"

	got, ok := st.Max()
	if !ok {
		t.Error("st.Max(): got no min value")
	}

	if got != want {
		t.Errorf("st.Max(): got unexpected value %v; want %v", got, want)
	}
}

func TestSearchTree_Max_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.Max()
	if ok {
		t.Errorf("st.Max(): got unexpected min value %v", got)
	}
}

func TestMultiValueSearchTree_Max(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3", "node3.1")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	want := []string{"node3", "node3.1"}

	got, ok := st.Max()
	if !ok {
		t.Error("st.Max(): got no min value")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("st.Max(): got unexpected value %v; want %v", got, want)
	}
}

func TestMultiValueSearchTree_Max_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	got, ok := st.Max()
	if ok {
		t.Errorf("st.Max(): got unexpected min value %v", got)
	}
}

func TestSearchTree_Ceil(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		start   int
		end     int
		wantOK  bool
		wantVal string
	}{
		{
			start:   9,
			end:     16,
			wantOK:  true,
			wantVal: "node5",
		},
		{
			start:   18,
			end:     20,
			wantOK:  true,
			wantVal: "node3",
		},
		{
			start:   7,
			end:     10,
			wantOK:  true,
			wantVal: "node6",
		},
		{
			start:  25,
			end:    30,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Ceil(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Ceil(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if got != tc.wantVal {
				t.Errorf("st.Ceil(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVal)
			}
		})
	}
}

func TestSearchTree_Ceil_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	start, end := 10, 15
	got, ok := st.Ceil(start, end)
	if ok {
		t.Errorf("st.Ceil(%v, %v): got unexpected ceiling value %v", start, end, got)
	}
}

func TestMultiValueSearchTree_Ceil(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3", "node3.1")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5", "node5.1")
	st.Insert(7, 10, "node6", "node6.1")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    9,
			end:      16,
			wantOK:   true,
			wantVals: []string{"node5", "node5.1"},
		},
		{
			start:    18,
			end:      20,
			wantOK:   true,
			wantVals: []string{"node3", "node3.1"},
		},
		{
			start:    7,
			end:      10,
			wantOK:   true,
			wantVals: []string{"node6", "node6.1"},
		},
		{
			start:  25,
			end:    30,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Ceil(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Ceil(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.Ceil(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestMultiValueSearchTree_Ceil_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	start, end := 10, 15
	got, ok := st.Ceil(start, end)
	if ok {
		t.Errorf("st.Ceil(%v, %v): got unexpected ceiling value %v", start, end, got)
	}
}

func TestSearchTree_Floor(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		start   int
		end     int
		wantOK  bool
		wantVal string
	}{
		{
			start:   9,
			end:     16,
			wantOK:  true,
			wantVal: "node6",
		},
		{
			start:   18,
			end:     20,
			wantOK:  true,
			wantVal: "node1",
		},
		{
			start:   7,
			end:     10,
			wantOK:  true,
			wantVal: "node6",
		},
		{
			start:  2,
			end:    4,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Floor(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Floor(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if got != tc.wantVal {
				t.Errorf("st.Floor(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVal)
			}
		})
	}
}

func TestSearchTree_Floor_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	start, end := 10, 15
	got, ok := st.Floor(start, end)
	if ok {
		t.Errorf("st.Floor(%v, %v): got unexpected floor value %v", start, end, got)
	}
}

func TestMultiValueSearchTree_Floor(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1", "node1.1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6", "node6.1")

	testCases := []struct {
		start    int
		end      int
		wantOK   bool
		wantVals []string
	}{
		{
			start:    9,
			end:      16,
			wantOK:   true,
			wantVals: []string{"node6", "node6.1"},
		},
		{
			start:    18,
			end:      20,
			wantOK:   true,
			wantVals: []string{"node1", "node1.1"},
		},
		{
			start:    7,
			end:      10,
			wantOK:   true,
			wantVals: []string{"node6", "node6.1"},
		},
		{
			start:  2,
			end:    4,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got, ok := st.Floor(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.Floor(%v, %v): got ok value %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.Floor(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.wantVals)
			}
		})
	}
}

func TestMultiValueSearchTree_Floor_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	start, end := 10, 15
	got, ok := st.Floor(start, end)
	if ok {
		t.Errorf("st.Floor(%v, %v): got unexpected floor value %v", start, end, got)
	}
}

func TestSearchTree_Rank(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		start int
		end   int
		want  int
	}{
		{
			start: 9,
			end:   16,
			want:  3,
		},
		{
			start: 18,
			end:   20,
			want:  5,
		},
		{
			start: 5,
			end:   8,
			want:  1,
		},
		{
			start: 2,
			end:   4,
			want:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got := st.Rank(tc.start, tc.end)
			if got != tc.want {
				t.Errorf("st.Rank(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.want)
			}
		})
	}
}

func TestMultiValueSearchTree_Rank(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		start int
		end   int
		want  int
	}{
		{
			start: 9,
			end:   16,
			want:  3,
		},
		{
			start: 18,
			end:   20,
			want:  5,
		},
		{
			start: 5,
			end:   8,
			want:  1,
		},
		{
			start: 2,
			end:   4,
			want:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			got := st.Rank(tc.start, tc.end)
			if got != tc.want {
				t.Errorf("st.Rank(%v, %v): got unexpected value %v; want %v", tc.start, tc.end, got, tc.want)
			}
		})
	}
}

func TestSearchTree_Select(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		k       int
		wantOK  bool
		wantVal string
	}{
		{
			k:       3,
			wantOK:  true,
			wantVal: "node5",
		},
		{
			k:      -1,
			wantOK: false,
		},
		{
			k:      8,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.k), func(t *testing.T) {
			got, ok := st.Select(tc.k)
			if ok != tc.wantOK {
				t.Errorf("st.Select(%v): got ok value %t; want %t", tc.k, ok, tc.wantOK)
			}

			if got != tc.wantVal {
				t.Errorf("st.Select(%v): got unexpected value %v; want %v", tc.k, got, tc.wantVal)
			}
		})
	}
}

func TestMultiValueSearchTree_Select(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	//defer mustBeValidTree(t, st)

	st.Insert(17, 19, "node1")
	st.Insert(5, 8, "node2")
	st.Insert(21, 24, "node3")
	st.Insert(4, 8, "node4")
	st.Insert(15, 18, "node5", "node5.1")
	st.Insert(7, 10, "node6")

	testCases := []struct {
		k        int
		wantOK   bool
		wantVals []string
	}{
		{
			k:        3,
			wantOK:   true,
			wantVals: []string{"node5", "node5.1"},
		},
		{
			k:      -1,
			wantOK: false,
		},
		{
			k:      8,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.k), func(t *testing.T) {
			got, ok := st.Select(tc.k)
			if ok != tc.wantOK {
				t.Errorf("st.Select(%v): got ok value %t; want %t", tc.k, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(got, tc.wantVals) {
				t.Errorf("st.Select(%v): got unexpected value %v; want %v", tc.k, got, tc.wantVals)
			}
		})
	}
}
