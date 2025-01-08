/*
Package sqliter provides utilities for querying SQL databases using iterators.

It defines a Querier interface for executing SQL queries and a Query function
that returns an iterator over the results of a query. The iterator yields
values of a specified type and any errors encountered during the iteration.
*/

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

// Query executes a query with the provided context and returns an iterator over the results.
// The scanner function is used to convert each row into a value of type T.
// If an error occurs during the query execution, the iterator will yield an error as the second return value.
func Query[T any](ctx context.Context, querier Querier, scanner func(row Row) (T, error), query string, args ...any) iter.Seq2[T, error] {
	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return func(yield func(T, error) bool) {
			var t T
			yield(t, err) // Yield an error if the query fails
		}
	}
	return func(yield func(T, error) bool) {
		defer rows.Close() // Ensure rows are closed after iteration
		for rows.Next() {
			v, err := scanner(rows) // Scan the row into a value of type T
			if !yield(v, err) {     // Yield the value and error
				return
			}
		}
		if err := rows.Err(); err != nil {
			var t T
			yield(t, err) // Yield an error if there was an issue during iteration
		}
	}
}
