package interval

import (
	"fmt"
	"testing"
)

func TestSearchTree_Delete(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(17, 19, 0)
	st.Insert(5, 8, 1)
	st.Insert(22, 24, 2)
	st.Insert(19, 23, 3)
	st.Insert(29, 35, 4)
	st.Insert(18, 20, 5)
	st.Insert(27, 28, 6)
	st.Insert(25, 28, 7)

	testCases := []struct {
		start int
		end   int
	}{
		{
			start: 27,
			end:   28,
		},
		{
			start: 17,
			end:   19,
		},
		{
			start: 5,
			end:   8,
		},
		{
			start: 18,
			end:   20,
		},
		{
			start: 29,
			end:   35,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			defer mustBeValidTree(t, st)

			if err := st.Delete(tc.start, tc.end); err != nil {
				t.Fatalf("st.Delete(%v, %v): got unexpected error %v", tc.start, tc.end, err)
			}

			got, ok := st.Find(tc.start, tc.end)
			if ok {
				t.Errorf("st.Find(%v, %v): got unexpected value %v", tc.start, tc.end, got)
			}
		})
	}
}

func TestSearchTree_Delete_EmptyTree(t *testing.T) {
	st := NewSearchTree[any](func(x, y int) int { return x - y })

	err := st.Delete(1, 10)
	if err != nil {
		t.Errorf("st.Delete(1, 10): got unexpected error %v", err)
	}
}

func TestSearchTree_Delete_NotFoundKey(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })
	st.Insert(20, 25, 0)

	err := st.Delete(20, 30)
	if err != nil {
		t.Errorf("st.Delete(20, 30): got unexpected error %v", err)
	}
}

func TestSearchTree_Delete_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewSearchTree[any](func(x, y int) int { return x - y })
		st.Insert(5, 10, nil)

		start, end := 10, 4
		err := st.Delete(start, end)
		if err == nil {
			t.Errorf("st.Delete(%v, %v): got nil error", start, end)
		}
	})
}

func TestSearchTree_DeleteMin(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(17, 19, 0)
	st.Insert(5, 8, 1)

	st.DeleteMin()

	if v, ok := st.Find(5, 8); ok {
		t.Errorf("Find(5, 8): got unexpected removed value: %v", v)
	}

	mustBeBalanced(t, st)

	st.DeleteMin()

	if v, ok := st.Find(17, 19); ok {
		t.Errorf("Find(17, 19): got unexpected removed value: %v", v)
	}

	mustBeBalanced(t, st)

	st.DeleteMin()

	if got, want := st.Size(), 0; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}

func TestSearchTree_DeleteMax(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(22, 25, 1)
	st.Insert(5, 7, 2)
	st.Insert(24, 26, 3)
	st.Insert(23, 25, 4)
	st.Insert(25, 27, 3)
	st.Insert(4, 10, 3)

	st.DeleteMax()

	if v, ok := st.Find(25, 27); ok {
		t.Errorf("Find(25, 27): got unexpected removed value: %v", v)
	}

	mustBeBalanced(t, st)

	st.DeleteMax()

	if v, ok := st.Find(24, 26); ok {
		t.Errorf("Find(24, 26): got unexpected removed value: %v", v)
	}

	mustBeBalanced(t, st)

	st.DeleteMax()

	if v, ok := st.Find(23, 25); ok {
		t.Errorf("Find(23, 25): got unexpected removed value: %v", v)
	}

	mustBeBalanced(t, st)

	st.DeleteMax()

	if got, want := st.Size(), 2; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}

func TestSearchTree_DeleteMax_EmptyTree(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	st.DeleteMax()

	if got, want := st.Size(), 0; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}

func TestMultiValueSearchTree_Delete(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(17, 19, 0)
	st.Insert(5, 8, 1)
	st.Insert(22, 24, 2)
	st.Insert(19, 23, 3)
	st.Insert(29, 35, 4)
	st.Insert(18, 20, 5)
	st.Insert(27, 28, 6)
	st.Insert(25, 28, 7)

	testCases := []struct {
		start int
		end   int
	}{
		{
			start: 27,
			end:   28,
		},
		{
			start: 17,
			end:   19,
		},
		{
			start: 5,
			end:   8,
		},
		{
			start: 18,
			end:   20,
		},
		{
			start: 29,
			end:   35,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.start, tc.end), func(t *testing.T) {
			//defer mustBeValidTree(t, st)

			if err := st.Delete(tc.start, tc.end); err != nil {
				t.Fatalf("st.Delete(%v, %v): got unexpected error %v", tc.start, tc.end, err)
			}

			got, ok := st.Find(tc.start, tc.end)
			if ok {
				t.Errorf("st.Find(%v, %v): got unexpected value %v", tc.start, tc.end, got)
			}
		})
	}
}

func TestMultiValueSearchTree_Delete_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })

	err := st.Delete(1, 10)
	if err != nil {
		t.Errorf("st.Delete(1, 10): got unexpected error %v", err)
	}
}

func TestMultiValueSearchTree_Delete_NotFoundKey(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })
	st.Insert(20, 25, 0)

	err := st.Delete(20, 30)
	if err != nil {
		t.Errorf("st.Delete(20, 30): got unexpected error %v", err)
	}
}

func TestMultiValueSearchTree_Delete_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewMultiValueSearchTree[any](func(x, y int) int { return x - y })
		st.Insert(5, 10, nil)

		start, end := 10, 4
		err := st.Delete(start, end)
		if err == nil {
			t.Errorf("st.Delete(%v, %v): got nil error", start, end)
		}
	})
}

func TestMultiValueSearchTree_DeleteMin(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(17, 19, 0)
	st.Insert(5, 8, 1)

	st.DeleteMin()

	if v, ok := st.Find(5, 8); ok {
		t.Errorf("st.Find(5, 8): got unexpected removed value: %v", v)
	}

	//mustBeBalanced(t, st)

	st.DeleteMin()

	if v, ok := st.Find(17, 19); ok {
		t.Errorf("st.Find(17, 19): got unexpected removed value: %v", v)
	}

	//mustBeBalanced(t, st)

	st.DeleteMin()

	if got, want := st.Size(), 0; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}

func TestMultiValueSearchTree_DeleteMax(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

	st.Insert(22, 25, 1)
	st.Insert(5, 7, 2)
	st.Insert(24, 26, 3)
	st.Insert(23, 25, 4)
	st.Insert(25, 27, 3)
	st.Insert(4, 10, 3)

	st.DeleteMax()

	if v, ok := st.Find(25, 27); ok {
		t.Errorf("Find(25, 27): got unexpected removed value: %v", v)
	}

	//mustBeBalanced(t, st)

	st.DeleteMax()

	if v, ok := st.Find(24, 26); ok {
		t.Errorf("Find(24, 26): got unexpected removed value: %v", v)
	}

	//mustBeBalanced(t, st)

	st.DeleteMax()

	if v, ok := st.Find(23, 25); ok {
		t.Errorf("Find(23, 25): got unexpected removed value: %v", v)
	}

	//mustBeBalanced(t, st)

	st.DeleteMax()

	if got, want := st.Size(), 2; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}

func TestMultiValueSearchTree_DeleteMax_EmptyTree(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

	st.DeleteMax()

	if got, want := st.Size(), 0; got != want {
		t.Errorf("st.Size(): got size %d; want %d", got, want)
	}
}
