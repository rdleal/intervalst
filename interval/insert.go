package interval

import (
	"fmt"
)

// Insert inserts the given val with the given start and end as the interval key.
// If there's already an interval key entry with the given start and end interval,
// it will be updated with the given val.
// Insert returns an error if the given end is less than or equal to the given start value.
func (st *SearchTree[V, T]) Insert(start, end T, val V) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	intervl := interval[V, T]{
		start: start,
		end:   end,
		val:   val,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	st.root = upsert(st.root, intervl, st.cmp)
	st.root.color = black

	return nil
}

func upsert[V, T any](h *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if h == nil {
		return newNode(intervl, red)
	}

	switch {
	case intervl.equal(h.interval.start, h.interval.end, cmp):
		h.interval = intervl
	case intervl.less(h.interval.start, h.interval.end, cmp):
		h.left = upsert(h.left, intervl, cmp)
	default:
		h.right = upsert(h.right, intervl, cmp)
	}

	if cmp.gt(intervl.end, h.maxEnd) {
		h.maxEnd = intervl.end
	}

	updateSize(h)

	return balanceNode(h, cmp)
}

type EmptyValueListError string

func newEmptyValueListError[V, T any](it interval[V, T], action string) error {
	s := fmt.Sprintf("multi value interval search tree: cannot %s empty value list for interval (%v, %v)", action, it.start, it.end)
	return EmptyValueListError(s)
}

func (e EmptyValueListError) Error() string {
	return string(e)
}

func (st *MultiValueSearchTree[V, T]) Insert(start, end T, vals ...V) error {
	intervl := interval[V, T]{
		start: start,
		end:   end,
		vals:  vals,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	if len(vals) == 0 {
		return newEmptyValueListError(intervl, "insert")
	}

	st.root = insert(st.root, intervl, st.cmp)
	st.root.color = black

	return nil
}

func insert[V, T any](h *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if h == nil {
		return newNode(intervl, red)
	}

	switch {
	case intervl.equal(h.interval.start, h.interval.end, cmp):
		h.interval.vals = append(h.interval.vals, intervl.vals...)
	case intervl.less(h.interval.start, h.interval.end, cmp):
		h.left = upsert(h.left, intervl, cmp)
	default:
		h.right = upsert(h.right, intervl, cmp)
	}

	if cmp.gt(intervl.end, h.maxEnd) {
		h.maxEnd = intervl.end
	}

	updateSize(h)

	return balanceNode(h, cmp)
}

func (st *MultiValueSearchTree[V, T]) Upsert(start, end T, vals ...V) error {
	intervl := interval[V, T]{
		start: start,
		end:   end,
		vals:  vals,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	if len(vals) == 0 {
		return newEmptyValueListError(intervl, "upsert")
	}

	st.root = upsert(st.root, intervl, st.cmp)
	st.root.color = black

	return nil
}
