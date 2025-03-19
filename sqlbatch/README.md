# SQL Batch

`sqlbatch` is a Go package designed for efficient batch processing of SQL queries. It provides a simple API to build SQL queries with placeholders and manage batches of arguments, reducing allocations and improving performance.

## Features

- Supports multiple placeholder formats (`$`, `:`, `?`).
- Caches query buffers and argument slices for reuse.
- Efficient batch processing with minimal allocations.

## Installation

```bash
go get github.com/Zamony/go/sqlbatch
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/Zamony/go/sqlbatch"
)

func main() {
	// Create a new batch with default options
	batch := sqlbatch.New(nil)
	defer batch.Close()

	// Append values to the batch
	batch.Append(1, "Alice", 25)
	batch.Append(2, "Bob", 30)

	// Build the SQL query
	query := batch.BuildQuery("INSERT INTO users (id, name, age) VALUES", "")

	// Get the arguments
	args := batch.BuildArguments()

	fmt.Println("Query:", query) // INSERT INTO users (id, name, age) VALUES ($1,$2,$3),($4,$5,$6);
	fmt.Println("Args:", args) // []any{1, "Alice", 25, 2, "Bob", 30}
}
```

## Benchmarks

```
BenchmarkBuildQueryStringsBuilder    39728750 ns/op	  2245079 B/op	   46991 allocs/op
BenchmarkBuildQueryBatch             15657675 ns/op	   906550 B/op	   20017 allocs/op
BenchmarkBuildQueryBatchCached        4725184 ns/op	   645697 B/op	   19825 allocs/op
```
