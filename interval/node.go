package interval

type color bool

const (
	red   color = true
	black color = false
)

type node[V, T any] struct {
	interval interval[V, T]
	maxEnd   T
	right    *node[V, T]
	left     *node[V, T]
	color    color
}

func newNode[V, T any](intervl interval[V, T], c color) *node[V, T] {
	return &node[V, T]{
		interval: intervl,
		maxEnd:   intervl.end,
		color:    c,
	}
}

func flipColors[T, V any](h *node[V, T]) {
	h.color = !h.color
	if h.left != nil {
		h.left.color = !h.left.color
	}
	if h.right != nil {
		h.right.color = !h.right.color
	}
}

func isRed[V, T any](h *node[V, T]) bool {
	if h == nil {
		return false
	}
	return h.color == red
}

func min[V, T any](h *node[V, T]) *node[V, T] {
	for h != nil && h.left != nil {
		h = h.left
	}
	return h
}

func max[V, T any](h *node[V, T]) *node[V, T] {
	for h != nil && h.right != nil {
		h = h.right
	}
	return h
}
