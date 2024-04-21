package interval

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"reflect"
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

func TestNewSearchTreeWithOptions_EmptyCmp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewSearchTreeWithOptions(nil): got execution without panic")
		}
	}()

	NewSearchTreeWithOptions[string, int](nil)
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

func TestMultiValueSearchTree_Height(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

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

func TestMultiValueSearchTree_Size(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

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

func TestMultiValueSearchTree_NilCmpFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewMultiValueSearchTree(nil): got execution without panic")
		}
	}()

	NewMultiValueSearchTree[string, int](nil)
}

func TestMultiValueSearchTreeWithOptions_NilCmpFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewMultiValueSearchTreeWithOptions(nil): got execution without panic")
		}
	}()

	NewMultiValueSearchTreeWithOptions[string, int](nil)
}

func TestMultiValueSearchTree_IsEmpty(t *testing.T) {
	st := NewMultiValueSearchTree[int](func(x, y int) int { return x - y })

	t.Run("EmptyTree", func(t *testing.T) {
		if got, want := st.IsEmpty(), true; got != want {
			t.Errorf("st.IsEmpty(): got unexpected value %t; want %t", got, want)
		}
	})

	t.Run("NotEmptyTree", func(t *testing.T) {
		st.Insert(10, 11, 0)

		if got, want := st.IsEmpty(), false; got != want {
			t.Errorf("st.IsEmpty(): got unexpected value %t; want %t", got, want)
		}
	})
}

func TestSearchTree_EncodingDecoding(t *testing.T) {
	tests := []struct {
		name string
		tree func() *SearchTree[string, int]
	}{
		{
			name: "with default options",
			tree: func() *SearchTree[string, int] {
				st := NewSearchTree[string, int](func(x, y int) int { return x - y })
				st.Insert(17, 19, "node1")
				st.Insert(5, 8, "node2")
				st.Insert(21, 24, "node3")
				st.Insert(21, 24, "node4")
				st.Insert(4, 4, "node5")

				return st
			},
		},
		{
			name: "with default options & empty",
			tree: func() *SearchTree[string, int] {
				return NewSearchTree[string, int](func(x, y int) int { return x - y })
			},
		},
		{
			name: "with interval point",
			tree: func() *SearchTree[string, int] {
				st := NewSearchTreeWithOptions[string, int](func(x, y int) int { return x - y }, TreeWithIntervalPoint())
				st.Insert(17, 19, "node1")
				st.Insert(5, 8, "node2")
				st.Insert(21, 24, "node3")
				st.Insert(21, 24, "node4")
				st.Insert(4, 4, "node5")

				return st
			},
		},
		{
			name: "with interval point & empty",
			tree: func() *SearchTree[string, int] {
				return NewSearchTreeWithOptions[string, int](func(x, y int) int { return x - y }, TreeWithIntervalPoint())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mustTestSearchTree_EncodingDecoding(t, tt.tree())
		})
	}
}

func mustTestSearchTree_EncodingDecoding(t *testing.T, st1 *SearchTree[string, int]) {
	t.Helper()

	st2 := NewSearchTree[string, int](func(x, y int) int { return x - y })

	defer mustBeValidTree(t, st1.root)
	defer mustBeValidTree(t, st2.root)

	b := mustEncodeTree(t, st1)
	mustDecodeTree(t, st2, b)

	// Roots should be equal
	if !reflect.DeepEqual(st1.root, st2.root) {
		t.Fatal("Roots are not equal")
	}

	// Configs should be equal: st2.config must be overridden by st1.config
	if !reflect.DeepEqual(st1.config, st2.config) {
		t.Fatal("Configs are not equal")
	}

	// After modifying the second tree,
	// roots should no longer be equal
	start, end := 2, 3

	err := st2.Insert(start, end, "node5")
	if err != nil {
		t.Fatalf("st.Insert(%v, %v): got unexpected error: %v", start, end, err)
	}

	if reflect.DeepEqual(st1.root, st2.root) {
		t.Fatal("Roots are still equal")
	}
}

func mustEncodeTree[V, T any](t *testing.T, st *SearchTree[V, T]) bytes.Buffer {
	t.Helper()
	var b bytes.Buffer

	w := bufio.NewWriter(&b)
	enc := gob.NewEncoder(w)

	err := enc.Encode(st)
	if err != nil {
		t.Fatalf("Encode: got unexpected error %v", err)
	}

	err = w.Flush()
	if err != nil {
		t.Fatalf("Flush: got unexpected error %v", err)
	}

	return b
}

func mustDecodeTree[V, T any](t *testing.T, st *SearchTree[V, T], b bytes.Buffer) {
	t.Helper()

	r := bufio.NewReader(&b)
	dec := gob.NewDecoder(r)

	err := dec.Decode(&st)
	if err != nil {
		t.Fatalf("Decode: got unexpected error %v", err)
	}
}

func TestMultiValueSearchTree_EncodingDecoding(t *testing.T) {
	tests := []struct {
		name string
		tree func() *MultiValueSearchTree[string, int]
	}{
		{
			name: "with default options",
			tree: func() *MultiValueSearchTree[string, int] {
				st := NewMultiValueSearchTree[string, int](func(x, y int) int { return x - y })
				st.Insert(17, 19, "node1")
				st.Insert(5, 8, "node2")
				st.Insert(21, 24, "node3")
				st.Insert(21, 24, "node4")
				st.Insert(4, 4, "node5")

				return st
			},
		},
		{
			name: "with default options & empty",
			tree: func() *MultiValueSearchTree[string, int] {
				return NewMultiValueSearchTree[string, int](func(x, y int) int { return x - y })
			},
		},
		{
			name: "with interval point",
			tree: func() *MultiValueSearchTree[string, int] {
				st := NewMultiValueSearchTreeWithOptions[string, int](func(x, y int) int { return x - y }, TreeWithIntervalPoint())
				st.Insert(17, 19, "node1")
				st.Insert(5, 8, "node2")
				st.Insert(21, 24, "node3")
				st.Insert(21, 24, "node4")
				st.Insert(4, 4, "node5")

				return st
			},
		},
		{
			name: "with interval point & empty",
			tree: func() *MultiValueSearchTree[string, int] {
				return NewMultiValueSearchTreeWithOptions[string, int](func(x, y int) int { return x - y }, TreeWithIntervalPoint())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mustTestMultiValueSearchTree_EncodingDecoding(t, tt.tree())
		})
	}
}

func mustTestMultiValueSearchTree_EncodingDecoding(t *testing.T, st1 *MultiValueSearchTree[string, int]) {
	t.Helper()

	st2 := NewMultiValueSearchTree[string, int](func(x, y int) int { return x - y })

	defer mustBeValidTree(t, st1.root)
	defer mustBeValidTree(t, st2.root)

	b := mustEncodeMultiValueTree(t, st1)
	mustDecodeMultiValueTree(t, st2, b)

	// Roots should be equal
	if !reflect.DeepEqual(st1.root, st2.root) {
		t.Fatal("Roots are not equal")
	}

	// Configs should be equal: st2.config must be overridden by st1.config
	if !reflect.DeepEqual(st1.config, st2.config) {
		t.Fatal("Configs are not equal")
	}

	// After modifying the second tree,
	// roots should no longer be equal
	start, end := 2, 3

	err := st2.Insert(start, end, "node5")
	if err != nil {
		t.Fatalf("st.Insert(%v, %v): got unexpected error: %v", start, end, err)
	}

	if reflect.DeepEqual(st1.root, st2.root) {
		t.Fatal("Roots are still equal")
	}
}

func mustEncodeMultiValueTree[V, T any](t *testing.T, st *MultiValueSearchTree[V, T]) bytes.Buffer {
	t.Helper()
	var b bytes.Buffer

	w := bufio.NewWriter(&b)
	enc := gob.NewEncoder(w)

	err := enc.Encode(st)
	if err != nil {
		t.Fatalf("Encode: got unexpected error %v", err)
	}

	err = w.Flush()
	if err != nil {
		t.Fatalf("Flush: got unexpected error %v", err)
	}

	return b
}

func mustDecodeMultiValueTree[V, T any](t *testing.T, st *MultiValueSearchTree[V, T], b bytes.Buffer) {
	t.Helper()

	r := bufio.NewReader(&b)
	dec := gob.NewDecoder(r)

	err := dec.Decode(&st)
	if err != nil {
		t.Fatalf("Decode: got unexpected error %v", err)
	}
}

func mustBeValidTree[V, T any](t *testing.T, root *node[V, T]) {
	t.Helper()

	mustBeBalanced(t, root)
	mustBeTwoThreeTree(t, root)
	mustHaveConsistentSize(t, root)
}

// Tests if all paths from root to leaf have the same number of blacks edges.
func mustBeBalanced[V, T any](t *testing.T, root *node[V, T]) {
	t.Helper()

	var black int
	for x := root; x != nil; x = x.Left {
		if !isRed(x) {
			black++
		}
	}

	if !isBalanced(root, black) {
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

	return isBalanced(h.Left, black) && isBalanced(h.Right, black)
}

// Tests if SearchTree is a 2-3 tree as left-leaning red black tree has a 1-1 correspondence to a 2-3 tree.
// For more on that, see https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
func mustBeTwoThreeTree[V, T any](t *testing.T, root *node[V, T]) {
	t.Helper()

	if !isTwoThreeTree(root) {
		t.Fatalf("Interval Tree is not a 2-3 tree")
	}
}

func isTwoThreeTree[V, T any](h *node[V, T]) bool {
	if h == nil {
		return true
	}

	if isRed(h.Right) {
		return false
	}

	if isRed(h.Left) && isRed(h.Right) {
		return false
	}

	return isTwoThreeTree(h.Left) && isTwoThreeTree(h.Right)
}

// Tests if the SearchTree nodes have consistent size.
func mustHaveConsistentSize[V, T any](t *testing.T, root *node[V, T]) {
	t.Helper()

	if !isSizeConsistent(root) {
		t.Fatalf("Interval Tree size is not consistent")
	}
}

func isSizeConsistent[V, T any](h *node[V, T]) bool {
	if h == nil {
		return true
	}

	if h.Size != size(h.Left)+size(h.Right)+1 {
		return false
	}

	return isSizeConsistent(h.Left) && isSizeConsistent(h.Right)
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
