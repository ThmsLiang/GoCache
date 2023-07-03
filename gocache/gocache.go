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
	"gocache/singleflight"
	pb "gocache/gocachepb"
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
	peers PeerPicker

	loader *singleflight.Group
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
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
		loader: &singleflight.Group{},
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
// update: add loader.Do to make sure load will only be executed once
func (g *Group) load(key string) (value ByteView, err error) {

	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GoCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	}) 
	
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key: key,
	}

	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
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



