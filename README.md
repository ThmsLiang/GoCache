# GoCache

GoCache is a distributed cache system implemented by Go. It is a simplified version of [groupcache](https://github.com/golang/groupcache), which is a go version of `memcached`.

## Credit

This project follows blog [GeeCache](https://geektutu.com/post/geecache.html).

## Features

GoCache supported the following features:

- Local cache and HTTP-based distributed cache.
- LRU (Least Recent Use) cache strategy.
- Avoid cache breakdown using `sync.Mutex` from Go.
- Choose nodes with consistent hashing for load balance.
- Optimize binary communication between two nodes with Protobuf.

## Structure

- `lru/lru.go`: Implemented Least-Recent-Use strategy for cache
- `byteview.go`: Abstruct and insulation of byte array in cache
- `cache.go`: Wrap lru and mutex for concurrency control
- `gocache.go`: Interact with outside, main procedure for get and update cache
