Minimal iterator-based SQL query library for Go.

## API

### QueryIter
```go
func QueryIter[T any](
    ctx context.Context,
    querier Querier,
    query string,
    args []any,
    scan RowScanner[T],
) iter.Seq2[T, error]
```

Streams results one by one. Stops early on error or when the caller breaks iteration.

### QuerySlice
```go
func QuerySlice[T any](
    ctx context.Context,
    querier Querier,
    query string,
    args []any,
    scan RowScanner[T],
) ([]T, error)
```

Returns all results as a slice. Returns early on first error.
