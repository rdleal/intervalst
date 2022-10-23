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
