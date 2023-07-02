Simple error handling primitives:
* stack traces;
* errors wrapping;
* compatible with stdlib;
* no external dependencies;

**API**
```go
// New creates an error with message and a stacktrace
err := errors.New("new")
err := errors.Newf("user %q doesn't exist", userID)

// Wrap adds a stacktrace to the error if it doesn't have one
err := errors.Wrap(err)
err := errors.Wrapf(err, "user %q doesn't exist", userID)

// Is reports whether any error in err's tree matches target
if errors.Is(err, ErrNotExists) {}
if errors.As(err, &valueErr) {}

// Combine multiple errors into one
err = errors.Combine(err, closeErr)

// Create constant errors
const ErrNotExists = errors.ConstantError("doesn't exist")
```
