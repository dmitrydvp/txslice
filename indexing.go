package txslice

import (
	"context"
	"sync"
)

type indexOpType int

const (
	maxIndexItems = 99_999

	indexOpAdd indexOpType = iota
	indexOpDelete
)

type txIndexInterface[T any] interface {
	Get(key any) (*T, bool)

	Add(v ...*T)
	Remove(v ...*T)

	Wait()
	Close()
}

type indexKeyFunc[T any, K comparable] func(v *T) K

type indexOperation[T any, K comparable] struct {
	typ   indexOpType
	value []*T
}

type index[T any, K comparable] struct {
	mu   sync.RWMutex
	data map[K]*T

	keyFn indexKeyFunc[T, K]

	ch chan indexOperation[T, K]
	wg sync.WaitGroup

	internalContext struct {
		ctx    context.Context
		cancel context.CancelFunc
		once   sync.Once
	}
}

func NewIndex[T any, K comparable](ctx context.Context, txsliceData *TxSlice[T], keyFn indexKeyFunc[T, K], buffSize int) {
	internalCtx, cancel := context.WithCancel(ctx)

	if buffSize == 0 {
		buffSize = 2048
	}

	idx := &index[T, K]{
		data:  make(map[K]*T, len(txsliceData.data)),
		keyFn: keyFn,
		ch:    make(chan indexOperation[T, K], buffSize),
		internalContext: struct {
			ctx    context.Context
			cancel context.CancelFunc
			once   sync.Once
		}{ctx: internalCtx, cancel: cancel},
	}

	for _, v := range txsliceData.data {
		idx.data[keyFn(v)] = v
	}

	go idx.loop()

	txsliceData.indexing = idx
}

func (i *index[T, K]) loop() {
	for {
		select {
		case <-i.internalContext.ctx.Done():
			return

		case op := <-i.ch:
			switch op.typ {
			case indexOpAdd:
				i.add(op.value)

			case indexOpDelete:
				i.delete(op.value)
			}

			i.wg.Done()
		}
	}
}

func (i *index[T, K]) Wait() {
	i.wg.Wait()
}

func (i *index[T, K]) add(values []*T) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, item := range values {
		i.data[i.keyFn(item)] = item
	}
}

func (i *index[T, K]) delete(values []*T) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, item := range values {
		delete(i.data, i.keyFn(item))
	}
}

func (i *index[T, K]) Add(v ...*T) {
	if i.internalContext.ctx.Err() != nil {
		return
	}

	isAsync := len(v) > maxIndexItems

	if isAsync {
		i.wg.Add(1)
		i.ch <- indexOperation[T, K]{typ: indexOpAdd, value: v}

		return
	}

	i.add(v)
}

func (i *index[T, K]) Remove(v ...*T) {
	if i.internalContext.ctx.Err() != nil {
		return
	}

	isAsync := len(v) > maxIndexItems

	if isAsync {
		i.wg.Add(1)
		i.ch <- indexOperation[T, K]{typ: indexOpDelete, value: v}

		return
	}

	i.delete(v)
}

func (i *index[T, K]) Get(key any) (*T, bool) {
	if i.internalContext.ctx.Err() != nil {
		return nil, false
	}

	k, ok := key.(K)
	if !ok {
		return nil, false
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	v, ok := i.data[k]

	return v, ok
}

func (i *index[T, K]) Close() {
	i.internalContext.once.Do(func() {
		i.internalContext.cancel()
		close(i.ch)
		i.data = nil
	})
}
