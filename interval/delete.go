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

	st.root = delete(st.root, intervl, st.cmp)
	if st.root != nil {
		st.root.color = black
	}

	return nil
}

func delete[V, T any](h *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if h == nil {
		return nil
	}

	if intervl.less(h.interval.start, h.interval.end, cmp) {
		if h.left != nil && !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h, cmp)
		}
		h.left = delete(h.left, intervl, cmp)
	} else {
		if isRed(h.left) {
			h = rotateRight(h, cmp)
		}
		if h.interval.equal(intervl.start, intervl.end, cmp) && h.right == nil {
			return nil
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h, cmp)
		}
		if h.interval.equal(intervl.start, intervl.end, cmp) {
			minNode := min(h.right)
			h.interval = minNode.interval
			h.right = deleteMin(h.right, cmp)
		} else {
			h.right = delete(h.right, intervl, cmp)
		}
	}

	updateSize(h)

	return fixUp(h, cmp)
}

func deleteMin[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	if h.left == nil {
		return nil
	}

	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h, cmp)
	}

	h.left = deleteMin(h.left, cmp)

	updateSize(h)

	return fixUp(h, cmp)
}

// DeleteMin removes the smallest interval key and its associated value from the tree.
func (st *SearchTree[V, T]) DeleteMin() {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.root == nil {
		return
	}

	st.root = deleteMin(st.root, st.cmp)
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

	st.root = deleteMax(st.root, st.cmp)
	if st.root != nil {
		st.root.color = black
	}
}

func deleteMax[V, T any](h *node[V, T], cmp CmpFunc[T]) *node[V, T] {
	if isRed(h.left) {
		h = rotateRight(h, cmp)
	}

	if h.right == nil {
		return nil
	}

	if !isRed(h.right) && !isRed(h.right.left) {
		h = moveRedRight(h, cmp)
	}

	h.right = deleteMax(h.right, cmp)

	updateSize(h)

	return fixUp(h, cmp)
}
