package par_test

import (
	"testing"

	"github.com/Zamony/go/par"
)

func TestMemoryPoolInit(t *testing.T) {
	t.Parallel()

	const length = 8
	pool := par.MemoryPool[[]byte]{}
	pool.Init(func() []byte {
		return make([]byte, length)
	})

	buf := pool.Get()
	equal(t, len(buf), length)

	pool.Put(buf)
	equal(t, len(pool.Get()), length)
}

func TestMemoryPoolNoInit(t *testing.T) {
	t.Parallel()

	pool := par.MemoryPool[[]byte]{}
	equal(t, pool.Get(), []byte(nil))

	pool.Put([]byte("a"))
	equal(t, pool.Get(), []byte("a"))
}
