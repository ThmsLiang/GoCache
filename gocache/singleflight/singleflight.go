package singleflight

import "sync"

// ongoing or past requests
type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

// main structure for singleflight to manage calls
type Group struct {
	mu sync.Mutex
	m map[string]*call
}

// for same key, fn will only execute once no matter how many Do() is called
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string] *call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // wait for request to complete
		return c.val, c.err 
	}

	c := new(call)
	c.wg.Add(1) // Add lock before executing fn
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()  // release lock after fn is executed

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err


}