# GoCache

GoCache is a distributive cache system implemented by Go. It is a simplified version of [groupcache](https://github.com/golang/groupcache), which is a go version of `memcached`.

## Credit

This project follows blog [GeeCache](https://geektutu.com/post/geecache.html).

## Features

GoCache supported the following features:

- Local cache and HTTP-based distributive cache.
- LRU (Least Recent Use) cache strategy.
- Avoid cache breakdown with lock mechanism from Go.
- Choose nodes with consistent hashing for load balance.
- Optimize binary communication between two nodes with Protobuf.
