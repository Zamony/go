package par

import "sync"

// A WaitGroup waits for a collection of goroutines to finish.
// A WaitGroup must not be copied after first use.
type WaitGroup struct {
	wg sync.WaitGroup
}

// Go calls the given function in a new goroutine.
func (g *WaitGroup) Go(fun func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fun()
	}()
}

// Wait blocks until all goroutines are completed.
func (g *WaitGroup) Wait() {
	g.wg.Wait()
}
