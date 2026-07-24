package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Zamony/go/middleware"
)

// requestLogger is a simple middleware for logging requests to stdout.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[logger] started: %s\n", r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("[logger] completed: %s\n", r.URL.Path)
	})
}

// recoverer is a middleware that recovers from panics and returns a 500 error.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[recoverer] caught panic: %v\n", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ExampleChain demonstrates composing multiple middlewares
// and applying them to a single handler.
func ExampleChain() {
	// 1. Create a chain: logger will be the outermost, recoverer the innermost.
	chain := middleware.Chain(requestLogger, recoverer)

	// 2. Create a base handler.
	// Note: We use fmt.Println to write to stdout so the Example test can verify the output order.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello from handler")
	})

	// 3. Wrap the handler.
	wrapped := chain(handler)

	// 4. Simulate an HTTP request for demonstration.
	req := httptest.NewRequest(http.MethodGet, "/test-chain", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	// Output:
	// [logger] started: /test-chain
	// hello from handler
	// [logger] completed: /test-chain
}

// apiHandler is a custom type that properly implements the http.Handler interface.
type apiHandler struct{}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Writing to stdout for the Example test verification.
	fmt.Println("api response")
}

// ExampleMount demonstrates registering an http.Handler (e.g., a struct)
// with automatic middleware application.
func ExampleMount() {
	mux := http.NewServeMux()

	// Create a helper function for registration.
	handle := middleware.Mount(mux, requestLogger)

	// Register the route using the custom struct that implements http.Handler.
	handle("/api/v1", apiHandler{})

	// Simulate a request.
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	// Output:
	// [logger] started: /api/v1
	// api response
	// [logger] completed: /api/v1
}

// ExampleMountFunc demonstrates the most common use case:
// registering a standard handler function with automatic middleware application,
// without needing to manually wrap it in http.HandlerFunc.
func ExampleMountFunc() {
	mux := http.NewServeMux()

	// Create a helper function for registering handler functions.
	handle := middleware.MountFunc(mux, requestLogger)

	// Register a route. Notice we pass a plain function, not an http.Handler.
	// MountFunc handles the conversion automatically.
	handle("/users", func(w http.ResponseWriter, r *http.Request) {
		// Writing to stdout for the Example test verification.
		fmt.Println("list of users")
	})

	// Simulate a request.
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	// Output:
	// [logger] started: /users
	// list of users
	// [logger] completed: /users
}
