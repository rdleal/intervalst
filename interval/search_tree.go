// Package interval provides a generic interval tree implementation.
//
// An interval tree is a data structure useful for storing values associated with intervals,
// and efficiently search those values based on intervals that overlap with any given interval.
// This generic implementation uses a self-balancing binary search tree algorithm, so searching
// for any intersection has a worst-case time-complexity guarantee of <= 2 log N, where N is the number of items in the tree.
//
// For more on interval trees, see https://en.wikipedia.org/wiki/Interval_tree
//
// To create a tree with time.Time as interval key type and string as value type:
//	cmpFn := func(t1, t2 time.Time) int {
//	  switch{
//	  case t1.After(t2): return 1
//	  case t1.Before(t2): return -1
//	  default: return 0
//	  }
//	}
// 	st := interval.NewSearchTree[string](cmpFn)
package interval

import (
	"math"
	"sync"
)

// SearchTree is a generic type representing the Interval Search Tree
// where V is a generic value type, and T is a generic interval key type.
type SearchTree[V, T any] struct {
	mu   sync.RWMutex // used to serialize read and write operations
	root *node[V, T]
	cmp  CmpFunc[T]
}

// NewSearchTree returns an initialized interval search tree.
// The cmp parameter is used for comparing total order of the interval key type T
// when inserting or looking up an interval in the tree.
// For more details on cmp, see the CmpFunc type.
// NewSearchTree will panic if cmp is nil.
func NewSearchTree[V, T any](cmp CmpFunc[T]) *SearchTree[V, T] {
	if cmp == nil {
		panic("NewSearchTree: comparison function cmp cannot be nil")
	}
	return &SearchTree[V, T]{
		cmp: cmp,
	}
}

// Height returns the max depth of the tree.
func (st *SearchTree[V, T]) Height() int {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return int(st.height(st.root))
}

func (st *SearchTree[V, T]) height(h *node[V, T]) float64 {
	if h == nil {
		return 0
	}

	return 1 + math.Max(st.height(h.left), st.height(h.right))
}

// Size returns the number of intervals in the tree.
func (st *SearchTree[V, T]) Size() int {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return size(st.root)
}

// IsEmpty returns true if the tree is empty; otherwise, false.
func (st *SearchTree[V, T]) IsEmpty() bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.root == nil
}
