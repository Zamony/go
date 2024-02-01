package par_test

import (
	"testing"

	"github.com/Zamony/go/par"
)

func TestAtomicValue(t *testing.T) {
	t.Parallel()

	number := par.AtomicValue[*int]{}
	_, ok := number.Load()
	equal(t, ok, false)

	one := ptrOfInt(1)
	_, ok = number.Swap(one)
	equal(t, ok, false)

	number.Store(one)
	value, ok := number.Load()
	equal(t, ok, true)
	equal(t, value, one)

	two := ptrOfInt(2)
	equal(t, number.CompareAndSwap(one, two), true)
	value, ok = number.Load()
	equal(t, ok, true)
	equal(t, value, two)

	value, ok = number.Swap(one)
	equal(t, ok, true)
	equal(t, value, two)
}

func ptrOfInt(v int) *int {
	return &v
}
