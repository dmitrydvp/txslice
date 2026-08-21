package txslice

import "fmt"

func (t *TxSlice[T]) BatchStart() *TxSlice[T] {
	b := New(t.data, Config{
		JournalCapacity:     t.journalCap,
		JournalCapacityStep: t.journalStep,
	})

	b.batchParent = t

	return b
}

func (t *TxSlice[T]) Batch(callback func(b *TxSlice[T]) error) error {
	b := t.BatchStart()

	if err := callback(b); err != nil {
		b.UndoBatch()
		return fmt.Errorf("callback: %w", err)
	}

	b.BatchAccept()

	return nil
}

func (t *TxSlice[T]) BatchAccept() {
	if t.batchParent == nil {
		return
	}

	t.IndexWait()

	t.batchParent.data = t.data
	t.batchParent.pushJournal(&operation[T]{
		typ:    batch,
		nested: t.journal,
	})

	t.Commit()

	t.batchParent = nil
}

func (t *TxSlice[T]) UndoBatch() {
	if t.batchParent == nil {
		return
	}

	t.batchParent = nil
}
