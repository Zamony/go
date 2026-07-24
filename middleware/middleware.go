package middleware

import "net/http"

// Chain returns a middleware that composes the given middlewares so that the
// first middleware in the list is the outermost (handles the request first)
// and the last middleware is the innermost (closest to the final handler).
//
// If no middlewares are provided, the returned function is an identity
// function that simply returns the original handler unchanged.
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// Mount returns a function that registers handlers on the given *http.ServeMux,
// automatically applying the provided middlewares (first is outermost).
//
// The returned function has the same signature as http.ServeMux.Handle, but it
// wraps each handler with the middleware chain before registering.
//
// Mount panics if mux is nil.
func Mount(mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) func(pattern string, handler http.Handler) {
	if mux == nil {
		panic("middleware: Mount called with nil ServeMux")
	}
	chain := Chain(middlewares...)
	return func(pattern string, handler http.Handler) {
		mux.Handle(pattern, chain(handler))
	}
}

// MountFunc is a convenience wrapper around Mount that accepts a handler
// function directly instead of an http.Handler, reducing boilerplate when
// the handler is defined as a bare function.
//
// MountFunc panics if mux is nil.
func MountFunc(mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mount := Mount(mux, middlewares...)
	return func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		mount(pattern, http.HandlerFunc(handler))
	}
}
