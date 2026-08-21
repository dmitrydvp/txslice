package txslice

import (
	"fmt"
)

type Config struct {
	IsAutoLatestSnap    bool
	JournalCapacity     int
	JournalCapacityStep int
}

type TxSlice[T any] struct {
	data        []*T
	journal     []*operation[T] // журнал изменений с точным описанием обратных действий
	journalCap  int
	journalStep int

	batchParent *TxSlice[T]
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
	}
}

func (t *TxSlice[T]) Commit() {
	res := make([]*T, t.Len())
	copy(res, t.data)

	t.data = res

	t.journal = newJournal[T](16)
}

func (t *TxSlice[T]) Rollback() {
	tempLog := "Operation index: %d; operation type: %d - %s; slice length: %d;\n"

	for i := len(t.journal) - 1; i >= 0; i-- {
		op := t.journal[i]
		indexFromBegin := ((len(t.journal) - 1) - i) + 1
		opMethod := ""

		switch op.typ {
		case opPush:
			opMethod = "push"
			t.undoPush(op)

		case opPop:
			opMethod = "pop"
			t.undoPop(op)

		case opShift:
			opMethod = "shift"
			t.undoShift(op)

		case opInsert:
			opMethod = "insert"
			t.undoInsert(op)

		case opSet:
			opMethod = "set"
			t.undoSet(op)

		case modSwap:
			opMethod = "mod swap"
			t.undoModSwap(op)

		case modMove:
			opMethod = "mod move"
			t.undoModMove(op)

		case batch:
			opMethod = "batch"
			t.batchRollback(op, indexFromBegin)

		}

		fmt.Printf(tempLog, indexFromBegin, op.typ, opMethod, t.Len())
	}

	t.journal = newJournal[T](16)
}

func (t *TxSlice[T]) batchRollback(op *operation[T], parentIndex int) {
	for i := len(op.nested) - 1; i >= 0; i-- {
		nestedOp := op.nested[i]
		indexFromBegin := ((len(op.nested) - 1) - i) + 1
		opMethod := ""

		switch nestedOp.typ {
		case opPush:
			opMethod = "push"
			t.undoPush(nestedOp)

		case opPop:
			opMethod = "pop"
			t.undoPop(nestedOp)

		case opShift:
			opMethod = "shift"
			t.undoShift(nestedOp)

		case opInsert:
			opMethod = "insert"
			t.undoInsert(nestedOp)

		case opSet:
			opMethod = "set"
			t.undoSet(nestedOp)

		case modSwap:
			opMethod = "mod swap"
			t.undoModSwap(nestedOp)

		case modMove:
			opMethod = "mod move"
			t.undoModMove(nestedOp)

		case batch:
			opMethod = "batch"
			t.batchRollback(nestedOp, indexFromBegin)

		}

		fmt.Printf("%d->\tOperation index: %d; operation type: %d - %s; slice length: %d;\n", parentIndex, indexFromBegin, nestedOp.typ, opMethod, t.Len())
	}

	t.journal = newJournal[T](16)
}
