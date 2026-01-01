// Package sqliter provides iterator for querying SQL databases.

package sqliter

import (
	"context"
	"database/sql"
	"iter"
)

// Querier is an interface for executing SQL queries.
type Querier interface {
	// QueryContext executes a query with the given context and arguments,
	// returning the result set as *sql.Rows.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// RowScanner scans the current row from *sql.Rows and returns a value of type T.
// The function is called for each row of the query result.
type RowScanner[T any] func(row *sql.Rows) (T, error)

// QueryIter executes a SQL query using the provided Querier interface and processes the results
// using the provided scanner function.
func QueryIter[T any](
	ctx context.Context,
	querier Querier,
	query string,
	args []any,
	scan RowScanner[T],
) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		rows, err := querier.QueryContext(ctx, query, args...)
		if err != nil {
			var t T
			yield(t, err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			t, err := scan(rows)
			if err != nil {
				yield(t, err)
				return
			}
			if !yield(t, nil) {
				return
			}
		}
		if err := rows.Err(); err != nil {
			var t T
			yield(t, err)
		}
	}
}

// QuerySlice executes an SQL query and returns all results as a slice.
// If an error occurs, it returns nil and an error.
func QuerySlice[T any](
	ctx context.Context,
	querier Querier,
	query string,
	args []any,
	scan RowScanner[T],
) ([]T, error) {
	elems := make([]T, 0, 32)
	for elem, err := range QueryIter(ctx, querier, query, args, scan) {
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}

	return elems, nil
}
