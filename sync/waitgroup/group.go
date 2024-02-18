package waitgroup

import "sync"

// A Group waits for a collection of goroutines to finish.
//
// A Group must not be copied after first use.
type Group struct {
	wg sync.WaitGroup
}

// Go calls the given function in a new goroutine.
func (g *Group) Go(fun func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fun()
	}()
}

// Wait blocks until all goroutines are completed.
func (g *Group) Wait() {
	g.wg.Wait()
}
