# GoCache

GoCache is a distributed cache system implemented by Go. It is a simplified version of [groupcache](https://github.com/golang/groupcache), which is a go version of `memcached`.

## Credit

This project follows blog [GeeCache](https://geektutu.com/post/geecache.html).

## Features

GoCache supported the following features:

- Local caching and HTTP-based distributed caching.
- LRU (Least Recent Use) cache strategy.
- Avoid cache breakdown using `sync.Mutex` from Go.
- Choose nodes with consistent hashing for load balance.
- Optimize binary communication between two peers with Protobuf.

## Structure

- Basic Data Structure:

  - `gocache/lru/lru.go`: Implemented Least-Recent-Use strategy for cache
  - `gocache/byteview.go`: Abstruct and insulation of byte array in cache
  - `gocache/cache.go`: Wrap lru and mutex for concurrency control
  - `gocache/gocache.go`: Interact with outside, main procedure to get and update cache

- HTTP related:

  - `gocache/http.go`: Define a http pool and parse http request for key search
  - `gocache/consistenthash/consistenthash.go`: Implemented consistent hasing
  - `gocache/peers.go`: Define `PeerPicker` and `PeerGetter` interface
  - `gocache/singleflight/singleflight.go`: Implemented `singleflight` to prevent cache breakdown

- Other:
  - `gocache/gocachepb/gocachepb.proto`: Use [Protobuf](https://github.com/protocolbuffers/protobuf) for peers to communicate
  - `run.sh`: start the server and run main.go
