Test helpers for use with the `google/go-cmp` package.

```
func Equal(t *testing.T, x, y any, opts ...cmp.Option)
func MustEqual(t *testing.T, x, y any, opts ...cmp.Option)
```

It may be handy to dot-import this package.
