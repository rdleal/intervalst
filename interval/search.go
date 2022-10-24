package interval

// Find returns the value which interval key exactly match the given start and end interval key.
// It returns true as the second return value if an exaclty matching interval key is found in the tree;
// otherwise, false.
func (st *SearchTree[V, T]) Find(start, end T) (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	var ok bool

	if st.root == nil {
		return val, ok
	}

	cur := st.root
	for cur != nil {
		switch {
		case cur.interval.equal(start, end, st.cmp):
			return cur.interval.val, true
		case cur.interval.less(start, end, st.cmp):
			cur = cur.right
		default:
			cur = cur.left
		}
	}

	return val, ok
}

// AnyIntersection returns a value which interval key intersects with the given start and end interval key.
// It returns true as the second return value if any intersection is found in the tree; otherwise, false.
func (st *SearchTree[V, T]) AnyIntersection(start, end T) (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	var ok bool

	if st.root == nil {
		return val, ok
	}

	cur := st.root
	for cur != nil {
		if cur.interval.intersects(start, end, st.cmp) {
			val, ok = cur.interval.val, true
			break
		}

		next := cur.left
		if cur.left == nil || st.cmp.gt(start, cur.left.maxEnd) {
			next = cur.right
		}

		cur = next
	}

	return val, ok
}

// AllIntersections returns a slice of values which interval key intersects with the given start and end interval key.
// It returns true as the second return value if any intersection is found in the tree; otherwise, false.
func (st *SearchTree[V, T]) AllIntersections(start, end T) ([]V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var vals []V
	if st.root == nil {
		return vals, false
	}

	st.searchInOrder(st.root, start, end, &vals)

	return vals, len(vals) > 0
}

func (st *SearchTree[V, T]) searchInOrder(h *node[V, T], start, end T, res *[]V) {
	if h.left != nil && st.cmp.gte(h.left.maxEnd, start) {
		st.searchInOrder(h.left, start, end, res)
	}

	if h.interval.intersects(start, end, st.cmp) {
		*res = append(*res, h.interval.val)
	}

	if h.right != nil && st.cmp.gte(h.right.maxEnd, start) {
		st.searchInOrder(h.right, start, end, res)
	}
}

// Min returns the value which interval key is the minimum interval in the tree.
// It returns false as the second return value if the tree is empty; otherwise, true.
func (st *SearchTree[V, T]) Min() (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	if st.root == nil {
		return val, false
	}

	val = min(st.root).interval.val

	return val, true
}

// Max returns the value which interval key is the maximum interval in the tree.
// It returns false as the second return value if the tree is empty; otherwise, true.
func (st *SearchTree[V, T]) Max() (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	if st.root == nil {
		return val, false
	}

	val = max(st.root).interval.val

	return val, true
}

// Ceil returns a value which interval key is the smallest interval greater than the given start and end interval.
// It returns true as the second return value if there's a ceiling interval key for the given start and end interval
// in the tree; otherwise, false.
func (st *SearchTree[V, T]) Ceil(start, end T) (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	if st.root == nil {
		return val, false
	}

	var ceil *node[V, T]

	cur := st.root
	for cur != nil {
		if cur.interval.equal(start, end, st.cmp) {
			return cur.interval.val, true
		}

		if cur.interval.less(start, end, st.cmp) {
			cur = cur.right
		} else {
			ceil = cur
			cur = cur.left
		}
	}

	if ceil == nil {
		return val, false
	}

	return ceil.interval.val, true
}

// Floor returns a value which interval key is the greatest interval lesser than the given start and end interval.
// It returns true as the second return value if there's a floor interval key for the given start and end interval
// in the tree; otherwise, false.
func (st *SearchTree[V, T]) Floor(start, end T) (V, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var val V
	if st.root == nil {
		return val, false
	}

	var floor *node[V, T]

	cur := st.root
	for cur != nil {
		if cur.interval.equal(start, end, st.cmp) {
			return cur.interval.val, true
		}

		if cur.interval.less(start, end, st.cmp) {
			floor = cur
			cur = cur.right
		} else {
			cur = cur.left
		}
	}

	if floor == nil {
		return val, false
	}

	return floor.interval.val, true
}
