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

func TestSearchTree_Time(t *testing.T) {
	t.Run("HasIntersection", func(t *testing.T) {
		st := NewSearchTree(timeCmp)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end)

		start, end = end.Add(1*time.Hour), end.Add(2*time.Hour)
		st.Insert(start, end)

		start, end = time.Now().Add(-(5 * time.Hour)), time.Now().Add(-(3 * time.Hour))
		st.Insert(start, end)

		wantStart, wantEnd := start, end

		start, end = start.Add(1*time.Hour), end.Add(1*time.Hour)

		gotStart, gotEnd, ok := st.AnyIntersection(start, end)
		if !ok {
			t.Errorf("st.AnyIntersection(%v, %v): got no intersection", start, end)
		}

		if !reflect.DeepEqual(gotStart, wantStart) {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected start value %v; want %v", start, end, gotStart, wantStart)
		}

		if !reflect.DeepEqual(gotEnd, wantEnd) {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected end value %v; want %v", start, end, gotEnd, wantEnd)
		}
	})

	t.Run("HasExactIntersection", func(t *testing.T) {
		st := NewSearchTree(timeCmp)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end)

		start, end = start.Add(2*time.Hour), end.Add(1*time.Hour)
		st.Insert(start, end)

		start, end = start.Add(-(5 * time.Hour)), end.Add(-(3 * time.Hour))
		st.Insert(start, end)

		wantStart, wantEnd := start, end

		gotStart, gotEnd, ok := st.AnyIntersection(start, end)
		if !ok {
			t.Errorf("st.AnyIntersection(%v, %v): got no intersection", start, end)
		}

		if !reflect.DeepEqual(gotStart, wantStart) {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected start value %v; want %v", start, end, gotStart, wantStart)
		}

		if !reflect.DeepEqual(gotEnd, wantEnd) {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected end value %v; want %v", start, end, gotEnd, wantEnd)
		}
	})

	t.Run("HasNoIntersection", func(t *testing.T) {
		st := NewSearchTree(timeCmp)

		start, end := time.Now(), time.Now().Add(1*time.Hour)
		st.Insert(start, end)

		start, end = start.Add(2*time.Hour), end.Add(1*time.Hour)
		st.Insert(start, end)

		start, end = start.Add(5*time.Hour), end.Add(3*time.Hour)
		st.Insert(start, end)

		start, end = start.Add(1*time.Hour), end.Add(1*time.Hour)

		gotStart, gotEnd, ok := st.AnyIntersection(start, end)
		if ok {
			t.Errorf("st.AnyIntersection(%v, %v): got unexpected intersection: %v - %v", start, end, gotStart, gotEnd)
		}
	})
}

func TestSearchTree_Insert_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewSearchTree(timeCmp)

		start, end := time.Now(), time.Now().Add(-(1 * time.Hour))
		err := st.Insert(start, end)
		if err == nil {
			t.Errorf("st.Insert(%v, %v): got nil error", start, end)
		}
	})
}

func TestSearchTree_Int(t *testing.T) {
	st := NewSearchTree(func(start, end int) int { return start - end })

	st.Insert(17, 19)
	st.Insert(5, 8)
	st.Insert(21, 24)
	st.Insert(4, 8)
	st.Insert(15, 18)
	st.Insert(7, 10)
	st.Insert(16, 22)

	testCases := []struct {
		start     int
		end       int
		wantStart int
		wantEnd   int
		wantOK    bool
	}{
		{
			start:     23,
			end:       25,
			wantOK:    true,
			wantStart: 21,
			wantEnd:   24,
		},
		{
			start:  12,
			end:    14,
			wantOK: false,
		},
		{
			start:     21,
			end:       23,
			wantOK:    true,
			wantStart: 16,
			wantEnd:   22,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			gotStart, gotEnd, ok := st.AnyIntersection(tc.start, tc.end)
			if ok != tc.wantOK {
				t.Errorf("st.AnyIntersection(%v, %v): got intersection %t; want %t", tc.start, tc.end, ok, tc.wantOK)
			}

			if !reflect.DeepEqual(gotStart, tc.wantStart) {
				t.Errorf("st.AnyIntersection(%v, %v): got unexpected start value %v; want %v", tc.start, tc.end, gotStart, tc.wantStart)
			}

			if !reflect.DeepEqual(gotEnd, tc.wantEnd) {
				t.Errorf("st.AnyIntersection(%v, %v): got unexpected end value %v; want %v", tc.start, tc.end, gotEnd, tc.wantEnd)
			}
		})
	}
}

func TestSearchTree_AnyIntersection_EmptyTree(t *testing.T) {
	st := NewSearchTree(func(i, j int) int { return j - i })

	gotStart, gotEnd, ok := st.AnyIntersection(1, 10)
	if ok {
		t.Errorf("st.AnyIntersect(1, 10): got unexpected intersection %v, %v", gotStart, gotEnd)
	}
}
