package txslice

func (t *TxSlice[T]) ModSwap(index1, index2 int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index1 < 0 || index2 < 0 || index1 >= t.Len() || index2 >= t.Len() {
		return
	}

	t.data[index1], t.data[index2] = t.data[index2], t.data[index1]

	t.pushJournal(&operation[T]{
		typ:     modSwap,
		indexes: []int{index2, index1},
		values:  []*T{t.data[index1], t.data[index2]},
	})
}

func (t *TxSlice[T]) ModMove(from, to int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if from < 0 || from >= t.Len() || to < 0 || to >= t.Len() {
		return
	}

	if from == to {
		return
	}

	val := t.data[from]
	copy(t.data[from:], t.data[from+1:])
	t.data = t.data[:len(t.data)-1]

	t.data = append(t.data, nil)
	copy(t.data[to+1:], t.data[to:])
	t.data[to] = val

	t.pushJournal(&operation[T]{
		typ:     modMove,
		indexes: []int{from, to},
		values:  []*T{val},
	})
}
