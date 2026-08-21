package txslice

func (t *TxSlice[T]) Slice() []*T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.data
}

func (t *TxSlice[T]) Len() int {
	return len(t.data)
}

func (t *TxSlice[T]) FirstElement() *T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.Len() == 0 {
		return nil
	}

	return t.data[0]
}

func (t *TxSlice[T]) LastElement() *T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.Len() == 0 {
		return nil
	}

	return t.data[t.Len()-1]
}

func (t *TxSlice[T]) MiddleElement() *T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.Len() == 0 {
		return nil
	}

	return t.data[(t.Len()-1)/2]
}

func (t *TxSlice[T]) At(index int) (*T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if index < 0 || index >= len(t.data) {
		return nil, false
	}

	return t.data[index], true
}

func (t *TxSlice[T]) IsEmpty() bool {
	return t.Len() == 0
}

func (t *TxSlice[T]) Min(keyFn func(*T) int) (*T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.data) == 0 {
		return nil, false
	}

	min := t.data[0]
	minKey := keyFn(min)

	for _, v := range t.data[1:] {
		k := keyFn(v)
		if k < minKey {
			min, minKey = v, k
		}
	}

	return min, true
}

func (t *TxSlice[T]) Max(keyFn func(*T) int) (*T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.data) == 0 {
		return nil, false
	}

	max := t.data[0]
	maxKey := keyFn(max)

	for _, v := range t.data[1:] {
		k := keyFn(v)
		if k > maxKey {
			max, maxKey = v, k
		}
	}

	return max, true
}

func (t *TxSlice[T]) Filter(predicate func(*T) bool) []*T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]*T, 0)

	for _, v := range t.data {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}

func (t *TxSlice[T]) Any(predicate func(*T) bool) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, v := range t.data {
		if predicate(v) {
			return true
		}
	}

	return false
}

func (t *TxSlice[T]) Find(predicate func(*T) bool) (int, *T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for i, v := range t.data {
		if predicate(v) {
			return i, v, true
		}
	}

	return -1, nil, false
}

func (t *TxSlice[T]) BinaryFind(key int, keyFn func(*T) int) (*T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	low, high := 0, len(t.data)-1

	for low <= high {
		mid := (low + high) / 2
		midVal := keyFn(t.data[mid])

		if midVal == key {
			return t.data[mid], true
		}

		if midVal < key {
			low = mid + 1
			continue
		}

		high = mid - 1

	}

	return nil, false
}
