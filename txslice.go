package txslice

import (
	"log"
	"sync"
)

type Config struct {
	IsAutoLatestSnap    bool
	JournalCapacity     int
	JournalCapacityStep int
	IsDebug             bool
}

type TxSlice[T any] struct {
	data        []*T
	journal     []*operation[T] // журнал изменений с точным описанием обратных действий
	journalCap  int
	journalStep int

	indexing txIndexInterface[T]

	batchParent *TxSlice[T]

	isDebug bool

	mu sync.RWMutex
}

func New[T any](data []*T, cfg Config) *TxSlice[T] {
	journalInitialCap := cfg.JournalCapacity

	if journalInitialCap == 0 {
		journalInitialCap = max(len(data)/4, defaultJournalMinCap)
	}

	journalCapStep := cfg.JournalCapacityStep

	if journalCapStep == 0 {
		journalCapStep = defaultJournalStep
	}

	return &TxSlice[T]{
		data:        data,
		journal:     newJournal[T](journalInitialCap),
		journalCap:  journalInitialCap,
		journalStep: journalCapStep,
		isDebug:     cfg.IsDebug,
	}
}

func (t *TxSlice[T]) IndexGet(key any) (*T, bool) {
	if t.indexing == nil {
		return nil, false
	}

	return t.indexing.Get(key)
}

func (t *TxSlice[T]) indexRemove(v ...*T) {
	if t.indexing == nil {
		return
	}

	t.indexing.Remove(v...)
}

func (t *TxSlice[T]) indexAdd(v ...*T) {
	if t.indexing == nil {
		return
	}

	t.indexing.Add(v...)
}

func (t *TxSlice[T]) IndexWait() {
	if t.indexing == nil {
		return
	}

	t.indexing.Wait()
}

func (t *TxSlice[T]) InTransaction(fn func(ts *TxSlice[T]) error) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := fn(t); err != nil {
		return err
	}

	t.Commit()

	return nil
}

func (t *TxSlice[T]) Commit() {
	t.mu.Lock()
	defer t.mu.Unlock()

	res := make([]*T, t.Len())
	copy(res, t.data)

	t.data = res

	t.journal = newJournal[T](16)
}

func (t *TxSlice[T]) Rollback() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := len(t.journal) - 1; i >= 0; i-- {
		op := t.journal[i]
		indexFromBegin := ((len(t.journal) - 1) - i) + 1

		opMethod := t.undoOperations(op, indexFromBegin)

		if t.isDebug {
			log.Printf("Operation index: %d; operation type: %d - %s; slice length: %d;\n", indexFromBegin, op.typ, opMethod, t.Len())
		}
	}

	t.journal = newJournal[T](16)
}

func (t *TxSlice[T]) batchRollback(op *operation[T], parentIndex int) {
	for i := len(op.nested) - 1; i >= 0; i-- {
		nestedOp := op.nested[i]
		indexFromBegin := ((len(op.nested) - 1) - i) + 1

		opMethod := t.undoOperations(nestedOp, parentIndex)

		if t.isDebug {
			log.Printf("%d-> Operation index: %d; operation type: %d - %s; slice length: %d;\n", parentIndex, indexFromBegin, nestedOp.typ, opMethod, t.Len())
		}
	}
}

func (t *TxSlice[T]) undoOperations(op *operation[T], parentIndex int) string {
	switch op.typ {
	case opPush:
		t.undoPush(op)

		return "push"

	case opPop:
		t.undoPop(op)

		return "pop"

	case opShift:
		t.undoShift(op)

		return "shift"

	case opInsert:
		t.undoInsert(op)

		return "insert"

	case opSet:
		t.undoSet(op)

		return "set"

	case modSwap:
		t.undoModSwap(op)

		return "mod swap"

	case modMove:
		t.undoModMove(op)

		return "mod move"

	case batch:
		t.batchRollback(op, parentIndex)

		return "batch"
	}

	return ""
}
