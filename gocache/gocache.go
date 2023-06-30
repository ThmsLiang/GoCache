// 									  Yes
// Receive key --> Check if in cache -----> return
//                 |  No                        		Yes
//                 |-----> Should get from other nodes  -----> communicate with other nodes --> return
//                             |  No
//                             |-----> call callback function, add to cache --> return

package gocache

import (
	"fmt"
	"log"
	"sync"
)

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name string
	getter Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cachBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("null Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cachBytes: cachBytes},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Check main cache first, if not then call load(key)
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is null")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[gocache] hit!")
		return v, nil
	}

	return g.load(key)
}

// either load value from local cache or outter cache
func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// use callback function to get value and add to maincache
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) 
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}



