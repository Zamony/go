[![Go Reference](https://pkg.go.dev/badge/github.com/Zamony/go/batcher.svg)](https://pkg.go.dev/github.com/Zamony/go/batcher)

Go library for batching items by count or time interval.

### Installation

```sh
go get github.com/Zamony/go/batcher
```

### Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Zamony/go/batcher"
)

func main() {
	// Create a new batcher that processes up to 100 items or every 1 second
	batcher := batcher.New[string](100, time.Second, func(ctx context.Context, items []string) {
		fmt.Printf("Processing batch of %d items: %v\n", len(items), items)
	})
	defer batcher.Close()

	// Add items to the batch
	ctx := context.Background()
	batcher.Add(ctx, "item1")
	batcher.Add(ctx, "item2", "item3")
	
	// Wait for processing
	time.Sleep(2 * time.Second)
}
```