package interval

// AnyIntersection returns an interval value that intersects with the given start and end interval key values.
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

// Find returns the exactly matching interval value for the given start and end interval key values.
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
