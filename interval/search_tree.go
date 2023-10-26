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
//
//	cmpFn := func(t1, t2 time.Time) int {
//	  switch{
//	  case t1.After(t2): return 1
//	  case t1.Before(t2): return -1
//	  default: return 0
//	  }
//	}
//	st := interval.NewSearchTree[string](cmpFn)
package interval

import (
	"sync"
)

type TreeConfig struct {
	allowIntervalPoint bool
}

// SearchTree is a generic type representing the Interval Search Tree
// where V is a generic value type, and T is a generic interval key type.
type SearchTree[V, T any] struct {
	mu     sync.RWMutex // used to serialize read and write operations
	root   *node[V, T]
	cmp    CmpFunc[T]
	config TreeConfig
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

type TreeOption func(*TreeConfig)

func TreeWithIntervalPoint() TreeOption {
	return func(c *TreeConfig) {
		c.allowIntervalPoint = true
	}
}

func NewSearchTreeWithOptions[V, T any](cmp CmpFunc[T], opts ...TreeOption) *SearchTree[V, T] {
	if cmp == nil {
		panic("NewSearchTreeWithOptions: comparison function cmp cannot be nil")
	}

	st := &SearchTree[V, T]{
		cmp: cmp,
	}

	for _, opt := range opts {
		opt(&st.config)
	}

	return st
}

// Height returns the max depth of the tree.
func (st *SearchTree[V, T]) Height() int {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return int(height(st.root))
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

// MultiValueSearchTree is a generic type representing the Interval Search Tree
// where V is a generic value type, and T is a generic interval key type.
// MultiValueSearchTree can store multiple values for a given interval key.
type MultiValueSearchTree[V, T any] SearchTree[V, T]

// NewMultiValueSearchTree returns an initialized multi value interval search tree.
// The cmp parameter is used for comparing total order of the interval key type T
// when inserting or looking up an interval in the tree.
// For more details on cmp, see the CmpFunc type.
// NewMultiValueSearchTree will panic if cmp is nil.
func NewMultiValueSearchTree[V, T any](cmp CmpFunc[T]) *MultiValueSearchTree[V, T] {
	if cmp == nil {
		panic("NewMultiValueSearchTree: comparison function cmp cannot be nil")
	}
	return &MultiValueSearchTree[V, T]{
		cmp: cmp,
	}
}

func NewMultiValueSearchTreeWithOptions[V, T any](cmp CmpFunc[T], opts ...TreeOption) *MultiValueSearchTree[V, T] {
	if cmp == nil {
		panic("NewMultiValueSearchTreeWithOptions: comparison function cmp cannot be nil")
	}

	st := &MultiValueSearchTree[V, T]{
		cmp: cmp,
	}

	for _, opt := range opts {
		opt(&st.config)
	}

	return st
}

// Height returns the max depth of the tree.
func (st *MultiValueSearchTree[V, T]) Height() int {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return int(height(st.root))
}

// Size returns the number of intervals in the tree.
func (st *MultiValueSearchTree[V, T]) Size() int {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return size(st.root)
}

// IsEmpty returns true if the tree is empty; otherwise, false.
func (st *MultiValueSearchTree[V, T]) IsEmpty() bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.root == nil
}
