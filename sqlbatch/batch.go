package sqlbatch

import (
	"bytes"
	"strconv"
	"sync"
)

// Constants defining the supported placeholder formats for SQL queries.
const (
	PlaceholderFormatDollar   byte = '$' // PostgreSQL-style placeholders ($1, $2, etc.)
	PlaceholderFormatColon    byte = ':' // Oracle-style placeholders (:1, :2, etc.)
	PlaceholderFormatQuestion byte = '?' // MySQL-style placeholders (?)
)

// Options holds configuration options for the SQL batch processing.
type Options struct {
	QueryCache        *queryCache     // Cache for query buffers to reduce allocations.
	ArgumentsCache    *argumentsCache // Cache for argument slices to reduce allocations.
	PlaceholderFormat byte            // Placeholder format to use in the generated SQL queries.
}

// defaultBatchOptions provides default configuration options.
var defaultBatchOptions = &Options{
	PlaceholderFormat: PlaceholderFormatDollar,
	QueryCache:        NewQueryCache(65536),
	ArgumentsCache:    NewArgumentsCache(65536),
}

// Batch represents a batch of SQL queries and their arguments.
type Batch struct {
	argumentsCache    *argumentsCache // Cache for argument slices.
	queryCache        *queryCache     // Cache for query buffers.
	buffer            *bytes.Buffer   // Buffer to build the SQL query.
	args              *[]any          // Accumulated arguments for the batch.
	batchSize         int             // Number of values per batch.
	placeholderFormat byte            // Placeholder format for the SQL query.
}

// New creates a new Batch instance with the provided options.
// Batch is not goroutine safe.
func New(options *Options) Batch {
	opts := parseOptions(options)
	return Batch{
		args:              opts.ArgumentsCache.get(),
		buffer:            opts.QueryCache.get(),
		placeholderFormat: opts.PlaceholderFormat,
		argumentsCache:    opts.ArgumentsCache,
		queryCache:        opts.QueryCache,
	}
}

// parseOptions merges provided options with defaults.
func parseOptions(opts *Options) *Options {
	if opts == nil {
		return defaultBatchOptions
	}
	o := *defaultBatchOptions
	if opts.PlaceholderFormat != 0 {
		o.PlaceholderFormat = opts.PlaceholderFormat
	}
	if opts.ArgumentsCache != nil {
		o.ArgumentsCache = opts.ArgumentsCache
	}
	if opts.QueryCache != nil {
		o.QueryCache = opts.QueryCache
	}
	return &o
}

// Append adds a set of values to the batch.
// Panics if the number of values differs from the batch size.
func (b *Batch) Append(values ...any) {
	if b.batchSize == 0 {
		b.batchSize = len(values)
	} else if len(values) != b.batchSize {
		panic("different sql batches sizes")
	}
	*b.args = append(*b.args, values...)
}

// BuildArguments returns the accumulated arguments for the batch.
func (b *Batch) BuildArguments() []any {
	return *b.args
}

// BuildQuery constructs the SQL query string using the provided prefix and suffix.
// Panics if batch size is zero.
func (b *Batch) BuildQuery(prefix, suffix string) string {
	if b.batchSize == 0 {
		panic("query batch size is zero")
	}

	b.buffer.WriteString(prefix)

	var placeholders []string
	if b.placeholderFormat != PlaceholderFormatQuestion {
		placeholders = getPlaceholders(len(*b.args))
	}
	count := 0
	batchCount := len(*b.args) / b.batchSize
	for range batchCount {
		if count > 0 {
			b.buffer.WriteByte(',')
		}
		b.buffer.WriteByte('(')
		lastIndex := b.batchSize - 1
		for i := range b.batchSize {
			count++
			b.buffer.WriteByte(b.placeholderFormat)
			if placeholders != nil {
				b.buffer.WriteString(placeholders[count])
			}
			if i != lastIndex {
				b.buffer.WriteByte(',')
			}
		}
		b.buffer.WriteByte(')')
	}

	if suffix != "" {
		b.buffer.WriteByte(' ')
		b.buffer.WriteString(suffix)
	}
	b.buffer.WriteByte(';')
	return b.buffer.String()
}

// Reset clears the batch for reuse, preserving allocated resources.
func (b *Batch) Reset() {
	b.batchSize = 0
	b.buffer.Reset()
	*b.args = (*b.args)[:0]
}

// Close releases resources used by the batch, returning them to the cache.
func (b *Batch) Close() {
	b.Reset()
	b.queryCache.put(b.buffer)
	b.argumentsCache.put(b.args)
}

var (
	placeholdersMutex = sync.Mutex{}  // Mutex to protect the placeholders slice.
	placeholders      = []string{"0"} // Precomputed placeholder strings.
)

// getPlaceholders returns the numeric suffix for a placeholder (e.g., "1" for "$1").
func getPlaceholders(count int) []string {
	placeholdersMutex.Lock()
	defer placeholdersMutex.Unlock()

	last := count + 1
	if count < len(placeholders) {
		values := placeholders[:last]
		return values
	}

	for i := len(placeholders); i <= count; i++ {
		placeholders = append(placeholders, strconv.Itoa(i))
	}

	return placeholders[:last]
}

// queryCache is a pool of bytes.Buffer instances to reduce allocations.
type queryCache struct {
	pool sync.Pool
}

// NewQueryCache creates a new query cache with the specified buffer size.
func NewQueryCache(size int) *queryCache {
	return &queryCache{sync.Pool{
		New: func() any {
			b := bytes.NewBuffer(nil)
			b.Grow(size)
			return b
		},
	}}
}

// get retrieves a buffer from the cache or creates a new one.
func (c *queryCache) get() *bytes.Buffer {
	return c.pool.Get().(*bytes.Buffer)
}

// put returns a buffer to the cache for reuse.
func (c *queryCache) put(buf *bytes.Buffer) {
	c.pool.Put(buf)
}

// argumentsCache is a pool of argument slices to reduce allocations.
type argumentsCache struct {
	pool sync.Pool
}

// NewArgumentsCache creates a new arguments cache with the specified slice capacity.
func NewArgumentsCache(size int) *argumentsCache {
	return &argumentsCache{sync.Pool{
		New: func() any {
			b := make([]any, 0, size)
			return &b
		},
	}}
}

// get retrieves an argument slice from the cache or creates a new one.
func (c *argumentsCache) get() *[]any {
	return c.pool.Get().(*[]any)
}

// put returns an argument slice to the cache for reuse.
func (c *argumentsCache) put(s *[]any) {
	c.pool.Put(s)
}
