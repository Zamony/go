package par_test

import (
	"reflect"
	"testing"

	"github.com/Zamony/go/par"
)

func TestMapSetGet(t *testing.T) {
	t.Parallel()

	var m par.Map[string, int]
	m.Set("a", 1)

	value, ok := m.SetIf("b", func(value int, exists bool) bool {
		return true
	}, func(int) int {
		return 2
	})
	equal(t, value, 2)
	equal(t, ok, true)

	value, ok = m.SetIf("b", func(value int, exists bool) bool {
		return !exists
	}, func(int) int {
		return 3
	})
	equal(t, value, 2)
	equal(t, ok, false)

	m.Set("c", 4)
	value, ok = m.SetIf("c", func(value int, exists bool) bool {
		return value == 4
	}, func(int) int {
		return 5
	})
	equal(t, value, 5)
	equal(t, ok, true)
	equal(t, m.Len(), 3)

	value, ok = m.Get("a")
	equal(t, ok, true)
	equal(t, value, 1)

	value, ok = m.Get("b")
	equal(t, ok, true)
	equal(t, value, 2)

	value, ok = m.Get("c")
	equal(t, ok, true)
	equal(t, value, 5)

	_, ok = m.Get("d")
	equal(t, ok, false)
}

func TestMapSetDelete(t *testing.T) {
	t.Parallel()

	var m par.Map[string, int]
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Delete("a")
	equal(t, m.DeleteIf("b", func(value int) bool {
		return value == 2
	}), true)
	equal(t, m.DeleteIf("c", func(value int) bool {
		return value == 100500
	}), false)
	equal(t, m.Len(), 1)

	_, ok := m.Get("a")
	equal(t, ok, false)

	_, ok = m.Get("b")
	equal(t, ok, false)

	value, ok := m.Get("c")
	equal(t, ok, true)
	equal(t, value, 3)
}

func TestMapForEach(t *testing.T) {
	t.Parallel()

	var m par.Map[string, int]
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	mit := map[string]int{}
	completed := m.ForEach(func(key string, value int) bool {
		mit[key] = value
		return true
	})
	equal(t, completed, true)
	equal(t, mit, map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	})

	m.Clear()
	equal(t, m.Len(), 0)
}

func TestMapConcurrent(t *testing.T) {
	t.Parallel()

	const key = "a"
	var m par.Map[string, int]
	var wg par.WaitGroup
	defer wg.Wait()

	for i := 1; i <= 10; i++ {
		i := i
		wg.Go(func() {
			m.Set(key, i)
		})
		wg.Go(func() {
			m.SetIf(key, func(int, bool) bool {
				return true
			}, func(v int) int {
				return v + 1
			})
		})
		wg.Go(func() {
			m.Get(key)
		})
		wg.Go(func() {
			m.Len()
		})
		wg.Go(func() {
			m.Delete(key)
		})
		wg.Go(func() {
			m.DeleteIf(key, func(int) bool {
				return true
			})
		})
		wg.Go(func() {
			m.ForEach(func(string, int) bool {
				return true
			})
		})
	}
}

func equal(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Not equal (-want, +got):\n- %+v\n+ %+v\n", want, got)
	}
}
