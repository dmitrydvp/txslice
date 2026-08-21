package txslice

func (t *TxSlice[T]) undoPush(op *operation[T]) {
	t.data = t.data[:t.Len()-op.countAppended]

	t.indexRemove(op.values...)
}

func (t *TxSlice[T]) undoPop(op *operation[T]) {
	t.data = append(t.data, op.values...)

	t.indexAdd(op.values...)
}

func (t *TxSlice[T]) undoShift(op *operation[T]) {
	t.data = append(op.values, t.data...)
	t.indexAdd(op.values...)
}

func (t *TxSlice[T]) undoInsert(op *operation[T]) {
	t.data = append(t.data[:op.indexes[0]], t.data[op.indexes[0]+1:]...)

	t.indexRemove(op.values...)
}

func (t *TxSlice[T]) undoSet(op *operation[T]) {
	t.indexRemove(t.data[op.indexes[0]])

	t.data[op.indexes[0]] = op.values[0]

	t.indexAdd(op.values...)
}

func (t *TxSlice[T]) undoModSwap(op *operation[T]) {
	t.data[op.indexes[0]], t.data[op.indexes[1]] = t.data[op.indexes[1]], t.data[op.indexes[0]]
}

func (t *TxSlice[T]) undoModMove(op *operation[T]) {
	from, to := op.indexes[0], op.indexes[1]

	if from >= 0 && from < t.Len() {
		val := t.data[from]
		t.data = append(t.data[:from], t.data[from+1:]...)

		if to >= t.Len() {
			t.data = append(t.data, val)
		} else {
			t.data = append(t.data[:to], append([]*T{val}, t.data[to:]...)...)
		}
	}
}
