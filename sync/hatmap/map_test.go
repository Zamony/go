package hatmap_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Zamony/go/sync/hatmap"
)

func TestMapSetGet(t *testing.T) {
	var m hatmap.Map[string, int]
	m.Set("a", 1)

	value, ok := m.SetIf("b", 2, func(currValue int) bool {
		return true
	})
	equal(t, value, 2)
	equal(t, ok, true)

	value, ok = m.SetIf("b", 3, func(currValue int) bool {
		return currValue == 0
	})
	equal(t, value, 2)
	equal(t, ok, false)

	m.Set("c", 4)
	value, ok = m.SetIf("c", 5, func(currValue int) bool {
		return currValue == 4
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
	m := hatmap.New[string, int](3)
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
	var m hatmap.Map[string, int]
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	mit := map[string]int{}
	for key, value := range m.All() {
		mit[key] = value
	}
	equal(t, mit, map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	})

	m.Clear()
	equal(t, m.Len(), 0)
}

func TestMapConcurrent(t *testing.T) {
	const key = "a"
	var m hatmap.Map[string, int]
	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 1; i <= 10; i++ {
		i := i
		goGroup(&wg, func() {
			m.Set(key, i)
		})
		goGroup(&wg, func() {
			m.SetIf(key, 77, func(currValue int) bool {
				return true
			})
		})
		goGroup(&wg, func() {
			m.Get(key)
		})
		goGroup(&wg, func() {
			m.Len()
		})
		goGroup(&wg, func() {
			m.Delete(key)
		})
		goGroup(&wg, func() {
			m.DeleteIf(key, func(int) bool {
				return true
			})
		})
		goGroup(&wg, func() {
			for k, v := range m.All() {
				_, _ = k, v
			}
		})
	}
}

func goGroup(wg *sync.WaitGroup, fun func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fun()
	}()
}

func equal(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Not equal (-want, +got):\n- %+v\n+ %+v\n", want, got)
	}
}
