package interval

import "fmt"

// InvalidIntervalError is a description of an invalid interval.
// Insert and Delete will return an InvalidIntervalError if the given interval key values are invalid.
type InvalidIntervalError string

// Error returns a string representation of the InvalidIntervalError error.
func (s InvalidIntervalError) Error() string {
	return string(s)
}

func newInvalidIntervalError[V, T any](it interval[V, T]) error {
	s := fmt.Sprintf("interval search tree invalid range: start value %v cannot be less than or equal to end value %v", it.start, it.end)
	return InvalidIntervalError(s)
}

// CmpFunc must return a nagative integer, zero or a positive interger as x is
// less than, equal to, or greater than y, respectively.
//
// CmpFunc imposes a total ordering on the given x and y values.
//
// It must also ensure that the relation is transitive: cmp(x, y) > 0 && cmp(y, z) > 0
// implies cmp(x, z) > 0.
type CmpFunc[T any] func(x, y T) int

func (f CmpFunc[T]) eq(x, y T) bool {
	return f(x, y) == 0
}

func (f CmpFunc[T]) lt(x, y T) bool {
	return f(x, y) < 0
}

func (f CmpFunc[T]) lte(x, y T) bool {
	return f(x, y) <= 0
}

func (f CmpFunc[T]) gt(x, y T) bool {
	return f(x, y) > 0
}

func (f CmpFunc[T]) gte(x, y T) bool {
	return f(x, y) >= 0
}

type interval[V, T any] struct {
	start T
	end   T
	val   V
}

func (it interval[V, T]) isInvalid(cmp CmpFunc[T]) bool {
	return cmp.lte(it.end, it.start)
}

func (it interval[V, T]) less(start, end T, cmp CmpFunc[T]) bool {
	return cmp.lt(it.start, start) || cmp.eq(it.start, start) && cmp.lt(it.end, end)
}

func (it interval[V, T]) intersects(start, end T, cmp CmpFunc[T]) bool {
	return cmp.lte(start, it.end) && cmp.lte(it.end, end)
}

func (it interval[V, T]) equal(start, end T, cmp CmpFunc[T]) bool {
	return cmp.eq(it.start, start) && cmp.eq(it.end, end)
}
