package interval

// Insert inserts the given val with the given start and end as the interval key.
// If there's already an interval key entry with the given start and end interval,
// it will be updated with the given val.
// Insert returns an error if the given end is less than or equal to the given start value.
func (st *SearchTree[V, T]) Insert(start, end T, val V) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	intervl := interval[V, T]{
		start: start,
		end:   end,
		val:   val,
	}

	if intervl.isInvalid(st.cmp) {
		return newInvalidIntervalError(intervl)
	}

	st.root = st.insert(st.root, intervl)
	st.root.color = black

	return nil
}

func (st *SearchTree[V, T]) insert(h *node[V, T], intervl interval[V, T]) *node[V, T] {
	if h == nil {
		return newNode(intervl, red, 1)
	}

	switch {
	case intervl.equal(h.interval.start, h.interval.end, st.cmp):
		h.interval = intervl
	case intervl.less(h.interval.start, h.interval.end, st.cmp):
		h.left = st.insert(h.left, intervl)
	default:
		h.right = st.insert(h.right, intervl)
	}

	if st.cmp.gt(intervl.end, h.maxEnd) {
		h.maxEnd = intervl.end
	}

	updateSize(h)

	return balanceNode(h, st.cmp)
}
