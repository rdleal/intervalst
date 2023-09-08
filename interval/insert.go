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

	st.root = upsert(st.root, intervl, st.cmp)
	st.root.color = black

	return nil
}

func upsert[V, T any](h *node[V, T], intervl interval[V, T], cmp CmpFunc[T]) *node[V, T] {
	if h == nil {
		return newNode(intervl, red)
	}

	switch {
	case intervl.equal(h.interval.start, h.interval.end, cmp):
		h.interval = intervl
	case intervl.less(h.interval.start, h.interval.end, cmp):
		h.left = upsert(h.left, intervl, cmp)
	default:
		h.right = upsert(h.right, intervl, cmp)
	}

	if cmp.gt(intervl.end, h.maxEnd) {
		h.maxEnd = intervl.end
	}

	updateSize(h)

	return balanceNode(h, cmp)
}
