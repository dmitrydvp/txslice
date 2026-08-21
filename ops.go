package txslice

func (t *TxSlice[T]) Push(n ...*T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.data = append(t.data, n...)

	t.pushJournal(&operation[T]{
		typ:           opPush,
		countAppended: len(n),
		values:        n,
	})

	t.indexAdd(n...)
}

func (t *TxSlice[T]) Pop() *T {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Len() == 0 {
		return nil
	}

	lastIndex := t.Len() - 1
	item := t.data[lastIndex]

	t.data = t.data[:lastIndex]

	t.pushJournal(&operation[T]{
		typ:     opPop,
		indexes: []int{lastIndex},
		values:  []*T{item},
	})

	t.indexRemove(item)

	return item
}

func (t *TxSlice[T]) Shift() *T {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Len() == 0 {
		return nil
	}

	item := t.data[0]

	t.data = t.data[1:]

	t.pushJournal(&operation[T]{
		typ:    opShift,
		values: []*T{item},
	})

	t.indexRemove(item)

	return item
}

func (t *TxSlice[T]) Insert(index int, item *T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index > t.Len() {
		return
	}

	if index == t.Len() {
		t.data = append(t.data, item)
	} else {
		t.data = append(t.data, nil)
		copy(t.data[index+1:], t.data[index:])
		t.data[index] = item
	}

	t.pushJournal(&operation[T]{
		typ:     opInsert,
		indexes: []int{index},
		values:  []*T{item},
	})

	t.indexAdd(item)
}

func (t *TxSlice[T]) Set(index int, item *T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if index < 0 || index >= t.Len() {
		return
	}

	t.indexRemove(t.data[index])

	t.data[index] = item

	t.pushJournal(&operation[T]{
		typ:     opSet,
		indexes: []int{index},
		values:  []*T{item},
	})

	t.indexAdd(item)
}
