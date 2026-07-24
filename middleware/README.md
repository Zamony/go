A minimal, zero-dependency library for composing HTTP middleware with Go's standard `net/http`.

## Quick Start

```go
package main

import (
	"net/http"
	"github.com/Zamony/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Create a mount function that applies logger and recoverer to all routes
	handle := middleware.MountFunc(mux, loggerMiddleware, recovererMiddleware)

	// Register routes cleanly, without repetitive chaining
	handle("GET /api/users", usersHandler)
	handle("POST /api/users", createUserHandler)

	http.ListenAndServe(":8080", mux)
}
```
