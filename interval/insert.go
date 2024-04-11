package interval

import (
	"fmt"
)

// Insert inserts the given val with the given start and end as the interval key.
// If there's already an interval key entry with the given start and end interval,
// it will be updated with the given val.
//
// Insert returns an InvalidIntervalError if the given end is less than or equal to the given start value.
func (st *SearchTree[V, T]) Insert(start, end T, val V) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	intervl := interval[V, T]{
		Start:      start,
		End:        end,
		Val:        val,
		AllowPoint: st.Config.AllowIntervalPoint,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	st.Root = upsert(st.Root, intervl, st.cmp)
	st.Root.Color = black

	return nil
}

func upsert[V, T any](n *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if n == nil {
		return newNode(intervl, red)
	}

	switch {
	case intervl.equal(n.Interval.Start, n.Interval.End, cmp):
		n.Interval = intervl
	case intervl.less(n.Interval.Start, n.Interval.End, cmp):
		n.Left = upsert(n.Left, intervl, cmp)
	default:
		n.Right = upsert(n.Right, intervl, cmp)
	}

	if cmp.gt(intervl.End, n.MaxEnd) {
		n.MaxEnd = intervl.End
	}

	updateSize(n)

	return balanceNode(n, cmp)
}

// EmptyValueListError is a description of an invalid list of values.
type EmptyValueListError string

// Error returns a string representation of the EmptyValueListError error.
func (e EmptyValueListError) Error() string {
	return string(e)
}

func newEmptyValueListError[V, T any](it interval[V, T], action string) error {
	s := fmt.Sprintf("multi value interval search tree: cannot %s empty value list for interval (%v, %v)", action, it.Start, it.End)
	return EmptyValueListError(s)
}

// Insert inserts the given vals with the given start and end as the interval key.
// If there's already an interval key entry with the given start and end interval,
// Insert will append the given vals to the exiting interval key.
//
// Insert returns an InvalidIntervalError if the given end is less than or equal to the given start value,
// or an EmptyValueListError if vals is an empty list.
func (st *MultiValueSearchTree[V, T]) Insert(start, end T, vals ...V) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	intervl := interval[V, T]{
		Start:      start,
		End:        end,
		Vals:       vals,
		AllowPoint: st.Config.AllowIntervalPoint,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	if len(vals) == 0 {
		return newEmptyValueListError(intervl, "insert")
	}

	st.Root = insert(st.Root, intervl, st.cmp)
	st.Root.Color = black

	return nil
}

func insert[V, T any](n *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if n == nil {
		return newNode(intervl, red)
	}

	switch {
	case intervl.equal(n.Interval.Start, n.Interval.End, cmp):
		n.Interval.Vals = append(n.Interval.Vals, intervl.Vals...)
	case intervl.less(n.Interval.Start, n.Interval.End, cmp):
		n.Left = insert(n.Left, intervl, cmp)
	default:
		n.Right = insert(n.Right, intervl, cmp)
	}

	if cmp.gt(intervl.End, n.MaxEnd) {
		n.MaxEnd = intervl.End
	}

	updateSize(n)

	return balanceNode(n, cmp)
}

// Upsert inserts the given vals with the given start and end as the interval key.
// If there's already an interval key entry with the given start and end interval,
// it will be updated with the given vals.
//
// Insert returns an InvalidIntervalError if the given end is less than or equal to the given start value,
// or an EmptyValueListError if vals is an empty list.
func (st *MultiValueSearchTree[V, T]) Upsert(start, end T, vals ...V) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	intervl := interval[V, T]{
		Start:      start,
		End:        end,
		Vals:       vals,
		AllowPoint: st.Config.AllowIntervalPoint,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	if len(vals) == 0 {
		return newEmptyValueListError(intervl, "upsert")
	}

	st.Root = upsert(st.Root, intervl, st.cmp)
	st.Root.Color = black

	return nil
}
