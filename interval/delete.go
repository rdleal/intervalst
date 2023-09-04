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
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h, st.cmp)
		}
		h.left = st.delete(h.left, intervl)
	} else {
		if isRed(h.left) {
			h = rotateRight(h, st.cmp)
		}
		if h.interval.equal(intervl.start, intervl.end, st.cmp) && h.right == nil {
			return nil
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h, st.cmp)
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

	return fixUp(h, st.cmp)
}

func (st *SearchTree[V, T]) deleteMin(h *node[V, T]) *node[V, T] {
	if h.left == nil {
		return nil
	}

	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h, st.cmp)
	}

	h.left = st.deleteMin(h.left)

	updateSize(h)

	return fixUp(h, st.cmp)
}

// DeleteMin removes the smallest interval key and its associated value from the tree.
func (st *SearchTree[V, T]) DeleteMin() {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.root == nil {
		return
	}

	st.root = st.deleteMin(st.root)
	if st.root != nil {
		st.root.color = black
	}
}

// DeleteMax removes the largest interval key and its associated value from the tree.
func (st *SearchTree[V, T]) DeleteMax() {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.root == nil {
		return
	}

	st.root = st.deleteMax(st.root)
	if st.root != nil {
		st.root.color = black
	}
}

func (st *SearchTree[V, T]) deleteMax(h *node[V, T]) *node[V, T] {
	if isRed(h.left) {
		h = rotateRight(h, st.cmp)
	}

	if h.right == nil {
		return nil
	}

	if !isRed(h.right) && !isRed(h.right.left) {
		h = moveRedRight(h, st.cmp)
	}

	h.right = st.deleteMax(h.right)

	updateSize(h)

	return fixUp(h, st.cmp)
}
