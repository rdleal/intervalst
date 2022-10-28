package interval

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewSearchTree_EmptyCmp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewSearchTree(nil): got execution without panic")
		}
	}()

	NewSearchTree[string, int](nil)
}

func TestSearchTree_Height(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	for i := 0; i < 255; i++ {
		err := st.Insert(i, i+1, i)
		if err != nil {
			t.Fatalf("Insert(%v, %v, %v): got unexpected error %v", i, i+1, i, err)
		}
	}

	if h := st.Height(); h%2 != 0 {
		t.Errorf("st.Height(): %v is not a power of 2", h)
	}
}

func TestSearchTree_Size(t *testing.T) {
	st := NewSearchTree[int](func(x, y int) int { return x - y })

	s := 20
	for i := 0; i < s; i++ {
		err := st.Insert(i, i+1, i)
		if err != nil {
			t.Fatalf("Insert(%v, %v, %v): got unexpected error %v", i, i+1, i, err)
		}
	}

	if got, want := st.Size(), s; got != want {
		t.Fatalf("st.Size(): got unexpected size %d; want %d", got, want)
	}

	err := st.Delete(4, 5)
	if err != nil {
		t.Fatalf("st.Delete(4, 5): Got unexpected error %v", err)
	}

	if got, want := st.Size(), s-1; got != want {
		t.Fatalf("st.Size(): got unexpected size %d; want %d", got, want)
	}

	err = st.Delete(15, 16)
	if err != nil {
		t.Fatalf("st.Delete(15, 16) Got unexpected error %v", err)
	}

	if got, want := st.Size(), s-2; got != want {
		t.Fatalf("st.Size(): got unexpected size %d; want %d", got, want)
	}

	// Already deleted, it should not affect the tree size.
	err = st.Delete(4, 5)
	if err != nil {
		t.Fatalf("st.Delete(4, 5): Got unexpected error %v", err)
	}

	if got, want := st.Size(), s-2; got != want {
		t.Fatalf("st.Size(): got unexpected size %d; want %d", got, want)
	}
}

func TestSearchTree_IsEmpty(t *testing.T) {
	t.Run("EmptyTree", func(t *testing.T) {
		st := NewSearchTree[int](func(x, y int) int { return x - y })

		if got, want := st.IsEmpty(), true; got != want {
			t.Errorf("st.IsEmpty(): got unexpected value %t; want %t", got, want)
		}
	})
	t.Run("NotEmptyTree", func(t *testing.T) {
		st := NewSearchTree[int](func(x, y int) int { return x - y })
		st.Insert(10, 11, 0)

		if got, want := st.IsEmpty(), false; got != want {
			t.Errorf("st.IsEmpty(): got unexpected value %t; want %t", got, want)
		}
	})
}

func mustBeValidTree[V, T any](t *testing.T, st *SearchTree[V, T]) {
	mustBeBalanced(t, st)
	mustBeTwoThreeTree(t, st)
}

// Tests if all paths from root to leaf have the same number of blacks edges.
func mustBeBalanced[V, T any](t *testing.T, st *SearchTree[V, T]) {
	var black int
	for x := st.root; x != nil; x = x.left {
		if !isRed(x) {
			black++
		}
	}

	if !isBalanced(st.root, black) {
		t.Fatal("Interval Tree is not balanced")
	}
}

func isBalanced[V, T any](h *node[V, T], black int) bool {
	if h == nil {
		return black == 0
	}
	if !isRed(h) {
		black--
	}

	return isBalanced(h.left, black) && isBalanced(h.right, black)
}

// Tests if SearchTree is a 2-3 tree as left-leaning red black tree has a 1-1 correspondence to a 2-3 tree.
func mustBeTwoThreeTree[V, T any](t *testing.T, st *SearchTree[V, T]) {
	if !isTwoThreeTree(st.root) {
		t.Fatalf("Interval Tree is not a 2-3 tree")
	}
}

func isTwoThreeTree[V, T any](h *node[V, T]) bool {
	if h == nil {
		return true
	}

	if isRed(h.right) {
		return false
	}

	if isRed(h.left) && isRed(h.right) {
		return false
	}

	return isTwoThreeTree(h.left) && isTwoThreeTree(h.right)
}

func testGenKeys(n int64) [][]int64 {
	rand.Seed(time.Now().UnixNano())
	res := make([][]int64, n)
	for i := 0; i < int(n); i++ {
		start := rand.Int63n(n)
		end := rand.Int63n(n-start+1) + start + 1
		res[i] = []int64{start, end}
	}

	return res

}

func BenchmarkSearchTree_Insert(b *testing.B) {
	testCases := []struct {
		n    int
		keys [][]int64
	}{
		{
			n:    10_000,
			keys: testGenKeys(10_000),
		},
		{
			n:    100_000,
			keys: testGenKeys(100_000),
		},
		{
			n:    1_000_000,
			keys: testGenKeys(1_000_000),
		},
		{
			n:    10_000_000,
			keys: testGenKeys(10_000_000),
		},
	}

	for _, tc := range testCases {
		b.Run(fmt.Sprint(tc.n), func(b *testing.B) {
			tree := NewSearchTree[int](func(x, y int64) int { return int(x - y) })

			for i := 0; i < b.N; i++ {
				for j, k := range tc.keys {
					err := tree.Insert(k[0], k[1], j)
					if err != nil {
						b.Fatalf("tree.Insert(%v, %v, %v): got unexpected error %v", k[0], k[1], i, err)
					}
				}
			}
		})
	}
}

func setupNewTree(keys int64) *SearchTree[int, int64] {
	rand.Seed(time.Now().UnixNano())
	st := NewSearchTree[int](func(x, y int64) int { return int(x - y) })

	for i, k := range testGenKeys(keys) {
		st.Insert(k[0], k[1], i)
	}

	return st
}

func testGenKey(n int64) (start, end int64) {
	rand.Seed(time.Now().UnixNano())
	start = rand.Int63n(n)
	end = rand.Int63n(n-start+1) + start + 1

	return
}

var result int

func BenchmarkSearchTree_AnyIntersection(b *testing.B) {
	testCases := []struct {
		keys int64
		tree *SearchTree[int, int64]
	}{
		{
			keys: 10_000,
			tree: setupNewTree(10_000),
		},
		{
			keys: 100_000,
			tree: setupNewTree(100_000),
		},
		{
			keys: 1_000_000,
			tree: setupNewTree(1_000_000),
		},
		{
			keys: 10_000_000,
			tree: setupNewTree(10_000_000),
		},
	}

	for _, tc := range testCases {
		start, end := testGenKey(tc.keys)
		b.Run(fmt.Sprint(tc.keys), func(b *testing.B) {
			var r int
			for i := 0; i < b.N; i++ {
				r, _ = tc.tree.AnyIntersection(start, end)
			}

			result = r
		})
	}
}

func BenchmarkSearchTree_Delete(b *testing.B) {
	testCases := []struct {
		keys int64
		tree *SearchTree[int, int64]
	}{
		{
			keys: 10_000,
			tree: setupNewTree(10_000),
		},
		{
			keys: 100_000,
			tree: setupNewTree(100_000),
		},
		{
			keys: 1_000_000,
			tree: setupNewTree(1_000_000),
		},
		{
			keys: 10_000_000,
			tree: setupNewTree(10_000_000),
		},
	}

	for _, tc := range testCases {
		start, end := testGenKey(tc.keys)
		b.Run(fmt.Sprint(tc.keys), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := tc.tree.Delete(start, end)
				if err != nil {
					b.Fatalf("tree.Delete(%v, %v, %v): got unexpected error %v", start, end, i, err)
				}
			}
		})
	}
}
