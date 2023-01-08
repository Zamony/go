Simple goroutine-safe rate limiter

Documentation: [https://pkg.go.dev/github.com/Zamony/golimit](https://pkg.go.dev/github.com/Zamony/golimit)

**Example**

```go
package main

import (
    "log"
    "time"

    "github.com/Zamony/go/golimit"
)

func main() {
	// Allow up to 10 calls per second
	lim := golimit.New(10, time.Second)

	for {
		if !lim.Limit(1) {
			log.Printf("ok")
		}
	}
}

```
