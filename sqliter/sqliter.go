// Package sqliter provides iterator for querying SQL databases.

package sqliter

import (
	"context"
	"database/sql"
	"iter"
)

// Row represents a single row returned by a query.
type Row interface {
	// Scan copies the columns from the current row into the values pointed to by dest.
	Scan(dest ...any) error
}

// Querier is an interface for executing SQL queries.
type Querier interface {
	// QueryContext executes a query with the given context and arguments,
	// returning the result set as *sql.Rows.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Query executes a SQL query using the provided Querier interface and processes the results
// using the provided scanner function.
func Query[T any](ctx context.Context, querier Querier, scanner func(row Row) (T, error), query string, args ...any) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		rows, err := querier.QueryContext(ctx, query, args...)
		if err != nil {
			var t T
			yield(t, err) // Yield an error if the query fails
			return
		}
		defer rows.Close() // Ensure rows are closed after iteration

		for rows.Next() {
			t, err := scanner(rows) // Scan the row into a value of type T
			if !yield(t, err) {     // Yield the value and error
				return
			}
		}
		if err := rows.Err(); err != nil {
			var t T
			yield(t, err) // Yield an error if there was an issue during iteration
		}
	}
}
