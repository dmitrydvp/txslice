# txslice

[![Go Reference](https://pkg.go.dev/badge/github.com/d1m4ek1/txslice.svg)](https://pkg.go.dev/github.com/d1m4ek1/txslice)
[![Go Report Card](https://goreportcard.com/badge/github.com/d1m4ek1/txslice)](https://goreportcard.com/report/github.com/d1m4ek1/txslice)

A **transactional slice** library for Go with undo/redo support, snapshots, batch operations, and optional indexing — all with minimal allocations and high performance.

⚠️ **Note:** This library is currently in active development. The API is not stable and may change.

## Status

🚧 **_Work in Progress_** – This library is under active development.

Expect breaking changes until version **v1.0.0.**

## Why txslice?

Sometimes you need **database-like transactions** but don’t want the overhead of a real DB.  
`txslice` provides an **in-memory transactional layer** for slices, making it easy to:

-   Apply changes and rollback safely.
-   Maintain historical states via snapshots.
-   Perform bulk operations with atomic commits.
-   Use indexing for O(1) lookups by arbitrary keys.
-   Build undo/redo systems or high-load pipelines.

## Features

-   **Transactions** — commit or rollback any change.
-   **Undo / Redo** — revert single or batch operations.
-   **Batch operations** — group multiple changes into one transaction.
-   **Snapshots with versioning** — save/restore slice states by tag.
-   **Indexing** — fast lookup by custom key extractor.
-   **Thread-safety** — built-in `sync.RWMutex` for concurrent access.
-   **Minimal allocations** — efficient memory usage even with millions of ops.

## Installation

```bash
go get github.com/d1m4ek1/txslice
```

## Example

### Simple Transactional Slice Operations

TxSlice Initialization and Simple Transactional Slice Operations

```golang
package main

import (
    "fmt"
    "github.com/d1m4ek1/txslice"
)

type Values struct {
    Val int
}

func main() {
    tx := txslice.New([]*Values{}, txslice.Config{})

	// Push elements
	tx.Push(&Values{1})
	tx.Push(&Values{2})

	fmt.Println("First element:", *tx.FirstElement()) // {1}

	// Create a snapshot
	tx.SetSnapshot("v1")

	// Commit example
	tx.Commit()

	// Transaction with rollback
	tx.Push(&Values{3})
	tx.Rollback() // rollback last transaction

	fmt.Println("Length after rollback:", tx.Len()) // 2
}
```

### Batch Operations

```go
package main

import (
	"fmt"

	"github.com/d1m4ek1/txslice"
)

type Values struct {
    Val int
}

func main() {
		// Initialize TxSlice
	tx := txslice.New([]*Values{{1}, {2}, {3}, {4}, {5}}, txslice.Config{})

	// Perform a batch of operations atomically
	if err := tx.Batch(func(b *txslice.TxSlice[Values]) error {
		b.Push(&Values{6}, &Values{7})

		b.ModSwap(0, 2)

		if b.FirstElement().Val != 3 {
			return errors.New("something error") // this error will be returned from Batch()
		}

		return nil
	}); err != nil {
		panic(err)
	}

	res := ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 3 2 1 4 5 6 7

	b := tx.BatchStart()

	b.ModMove(0, 3)

	b.BatchAccept()

	res = ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 2 1 4 3 5 6 7

	b = tx.BatchStart()

	b.ModMove(0, 3)

	b.UndoBatch()

	res = ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 2 1 4 3 5 6 7
}
```

## Performance

txslice is designed for high-load scenarios:

-   Minimal copying (lazy snapshots, journal-based rollback).

-   Optimized for millions of iterations.

-   Especially useful in route optimization, in-memory caching, undo/redo stacks.

**_Benchmarks coming soon._**

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.
Make sure to add/update tests as needed.
