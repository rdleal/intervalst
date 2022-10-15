package interval

import (
	"fmt"
)

type node[T any] struct {
	start  T
	end    T
	maxEnd T
	right  *node[T]
	left   *node[T]
}

func newNode[T any](start, end T) *node[T] {
	return &node[T]{
		start:  start,
		end:    end,
		maxEnd: end,
	}
}

type SearchTree[T any] struct {
	root *node[T]
	cmp  func(start, end T) int
}

func NewSearchTree[T any](cmp func(start, end T) int) *SearchTree[T] {
	return &SearchTree[T]{
		cmp: cmp,
	}
}

func (st *SearchTree[T]) Insert(start, end T) error {
	if st.cmp(end, start) <= 0 {
		return fmt.Errorf("interval search tree invalid range: start value %v cannot be less than or equal end value %v", start, end)
	}

	if st.root == nil {
		st.root = newNode(start, end)
		return nil
	}

	cur := st.root

	for cur != nil {
		if st.cmp(end, cur.maxEnd) > 0 {
			cur.maxEnd = end
		}

		if st.cmp(start, cur.start) < 0 || st.cmp(start, cur.start) == 0 && st.cmp(end, cur.end) < 0 {
			if cur.left == nil {
				cur.left = newNode(start, end)
				break
			}
			cur = cur.left
		} else {
			if cur.right == nil {
				cur.right = newNode(start, end)
				break
			}
			cur = cur.right
		}
	}

	return nil
}

func (st *SearchTree[T]) AnyIntersection(start, end T) (intersectStart, intersectEnd T, ok bool) {
	if st.root == nil {
		return
	}

	cur := st.root

	for cur != nil {
		if st.intersects(cur, start, end) {
			return cur.start, cur.end, true
		}

		next := cur.left
		if cur.left == nil || st.cmp(start, cur.left.maxEnd) > 0 {
			next = cur.right
		}

		cur = next
	}

	return
}

func (st *SearchTree[T]) intersects(n *node[T], start, end T) bool {
	return st.cmp(n.start, end) <= 0 && st.cmp(start, n.end) <= 0
}
