package lru

import "container/list"

type Cache struct {
	maxBytes int64                            // max number of bytes allowed for cache
	nBytes int64							  // number of bytes used
	ll *list.List					  // linkedlist to keep track of lru
	cache map[string]*list.Element			  // map key to a node in linkedlist
	OnEvicted func (key string, value Value)  // callback function when a node is evicted
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, OnEvicted func(string, Value)) *Cache{
	return &Cache {
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// Remove the least recent used kv
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if(c.OnEvicted != nil) {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add a new key-value or update value
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int{
	return c.ll.Len()
} 