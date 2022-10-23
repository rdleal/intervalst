package interval

// Find returns the value which interval key exactly match the given start and end interval key.
// It returns false as the second return value if no matching interval key is found in the tree.
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
// It returns false as the second return value if no intersection is found in the tree.
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
// It returns false as the second return value if no intersection is found in the tree.
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
