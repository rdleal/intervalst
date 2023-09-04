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
	size     int
}

func newNode[V, T any](intervl interval[V, T], c color, sz int) *node[V, T] {
	return &node[V, T]{
		interval: intervl,
		maxEnd:   intervl.end,
		color:    c,
		size:     sz,
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

func updateSize[V, T any](h *node[V, T]) {
	h.size = 1 + size(h.left) + size(h.right)
}

func size[V, T any](h *node[V, T]) int {
	if h == nil {
		return 0
	}
	return h.size
}

func updateMaxEnd[V, T any](h *node[V, T], cmp CmpFunc[T]) {
	h.maxEnd = h.interval.end
	if h.left != nil && cmp.gt(h.left.maxEnd, h.maxEnd) {
		h.maxEnd = h.left.maxEnd
	}

	if h.right != nil && cmp.gt(h.right.maxEnd, h.maxEnd) {
		h.maxEnd = h.right.maxEnd
	}
}

func rotateLeft[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	x.maxEnd = h.maxEnd
	h.color = red
	x.size = h.size

	updateSize(h)
	updateMaxEnd(h, cmp)
	return x
}

func rotateRight[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	x.maxEnd = h.maxEnd
	h.color = red
	x.size = h.size

	updateSize(h)
	updateMaxEnd(h, cmp)
	return x
}

func balanceNode[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h, cmp)
	}

	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h, cmp)
	}

	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}

	return h
}

func moveRedLeft[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	flipColors(h)
	if h.right != nil && isRed(h.right.left) {
		h.right = rotateRight(h.right, cmp)
		h = rotateLeft(h, cmp)
		flipColors(h)
	}
	return h
}

func moveRedRight[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	flipColors(h)
	if h.left != nil && isRed(h.left.left) {
		h = rotateRight(h, cmp)
		flipColors(h)
	}
	return h
}

func fixUp[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	updateMaxEnd(h, cmp)

	return balanceNode(h, cmp)
}
