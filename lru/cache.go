package lru

import "sync"

type Cache[K comparable, V any] struct {
	size int
	dict map[K]*Node[K, V]
	list List[K, V]
	seal sync.Mutex
}

type InvalidCacheSize struct{}

func (*InvalidCacheSize) Error() string {
	return "Invalid cache size (must be greater than 0)"
}

func New[K comparable, V any](size int) (*Cache[K, V], error) {
	if size < 1 {
		return nil, &InvalidCacheSize{}
	}

	cache := &Cache[K, V]{
		size: size,
		dict: make(map[K]*Node[K, V], size),
	}

	return cache, nil
}

func (c *Cache[K, V]) Cap() int {
	return c.size
}

func (c *Cache[K, V]) Len() int {
	return c.list.len
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.seal.Lock()
	defer c.seal.Unlock()

	if node, ok := c.dict[key]; ok {
		node.Value = value
		c.list.move(node)

		return
	}

	if c.Len() == c.Cap() {
		tail := c.list.head.prev

		c.list.remove(tail)
		delete(c.dict, tail.Key)
	}

	node := &Node[K, V]{Key: key, Value: value}

	c.list.insert(node)
	c.dict[key] = node
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.seal.Lock()
	defer c.seal.Unlock()

	if node, ok := c.dict[key]; ok {
		c.list.move(node)
		return node.Value, true
	} else {
		var empty V
		return empty, false
	}
}

func (c *Cache[K, V]) Del(key K) bool {
	c.seal.Lock()
	defer c.seal.Unlock()

	if node, ok := c.dict[key]; ok {
		c.list.remove(node)
		delete(c.dict, key)

		return true
	} else {
		return false
	}
}

// the following functions are subject to race condition

func (c *Cache[K, V]) GetOrPut(key K, value V) V {
	if v, ok := c.Get(key); ok {
		return v
	} else {
		c.Put(key, value)
		return value
	}
}

func (c *Cache[K, V]) GetOrPutFunc(key K, valueFunc func() V) V {
	if v, ok := c.Get(key); ok {
		return v
	} else {
		v := valueFunc()
		c.Put(key, v)

		return v
	}
}
