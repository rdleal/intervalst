package interval

// Delete removes the value associated with the given start and end interval key.
// Delete returns an error if the given end is less than or equal to the given start value.
func (st *SearchTree[V, T]) Delete(start, end T) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.root == nil {
		return nil
	}

	intervl := interval[V, T]{
		start: start,
		end:   end,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	st.root = st.delete(st.root, intervl)
	if st.root != nil {
		st.root.color = black
	}

	return nil
}

func (st *SearchTree[V, T]) delete(h *node[V, T], intervl interval[V, T]) *node[V, T] {
	if h == nil {
		return nil
	}

	if intervl.less(h.interval.start, h.interval.end, st.cmp) {
		if isRed(h.left) && !isRed(h.left.left) {
			h = st.moveRedLeft(h)
		}
		h.left = st.delete(h.left, intervl)
	} else {
		if isRed(h.left) {
			h = st.rotateRight(h)
		}
		if h.interval.equal(intervl.start, intervl.end, st.cmp) && h.right == nil {
			return nil
		}
		if !isRed(h.right) && h.right != nil && !isRed(h.right.left) {
			h = st.moveRedRight(h)
		}
		if h.interval.equal(intervl.start, intervl.end, st.cmp) {
			minNode := min(h.right)
			h.interval = minNode.interval
			h.right = st.deleteMin(h.right)
		} else {
			h.right = st.delete(h.right, intervl)
		}
	}

	updateSize(h)

	return st.fixUp(h)
}

func (st *SearchTree[V, T]) deleteMin(h *node[V, T]) *node[V, T] {
	if h.left == nil {
		return nil
	}

	if !isRed(h.left) && !isRed(h.left.left) {
		h = st.moveRedLeft(h)
	}

	h.left = st.deleteMin(h.left)

	updateSize(h)

	return st.fixUp(h)
}
