# lru-cache-go

this project is just a playground for learning Go, not intended for use in production (at least not yet)

## Example

```go
package main

import (
	"fmt"

	"github.com/netrusov/lru-cache-go/lru"
)

func main() {
	// make cache with 3 max keys and specify key/value shape
	cache, err := lru.New[int, int](3)

	// an error can be returned if the size is less than 1
	if err != nil {
		panic(err)
	}

	// the FIFO method is used to manage nodes
	// new nodes are always placed at the beginning
	cache.Put(1, 1) // state: [{1, 1}]
	cache.Put(2, 2) // state: [{2, 2}, {1, 1}]
	cache.Put(3, 3) // state: [{3, 3}, {2, 2}, {1, 1}]

	// if a node has been accessed, it will be moved to the beginning
	cache.Get(2) // state: [{2, 2}, {3, 3}, {1, 1}]

	// to distinguish a missing entry from a zero value, use "comma ok"
	if v, ok := cache.Get(2); ok {
		fmt.Println("entry found, value:", v)
	}

	// {1, 1} will be evicted due to size policy
	cache.Put(4, 4) // state: [{4, 4}, {2, 2}, {3, 3}]

	// {2, 2} will be updated with the new value and moved to the beginning
	cache.Put(2, 222) // state: [{2, 222}, {4, 4}, {3, 3}]

	if _, ok := cache.Get(1); !ok {
		fmt.Println("entry not found") // will be executed because "1" key is no longer present
	}
	// cache stays the same
	// state: [{2, 222}, {4, 4}, {3, 3}]

	cache.Del(4) // state: [{2, 222}, {3, 3}]
}
```
