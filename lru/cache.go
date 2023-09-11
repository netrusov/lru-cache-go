package lru

/*
TODO:
	- support concurrent access
	- debug weird behavior when there's only one element in the cache
	- possibly allow `nil` as a cache's head (may fix the bug above), although it will introduce a lot of nil-checks
*/

type Cache[K comparable, V any] struct {
	size int
	dict map[K]*Node[K, V]
	list List[K, V]
}

type InvalidCacheSize struct{}

func (*InvalidCacheSize) Error() string {
	return "Invalid cache size (must be greater than 1)"
}

func New[K comparable, V any](size int) (*Cache[K, V], error) {
	if size <= 1 {
		return nil, &InvalidCacheSize{}
	}

	cache := &Cache[K, V]{
		size: size,
		dict: make(map[K]*Node[K, V]),
	}

	return cache, nil
}

func (c *Cache[K, V]) evict() {
	if c.Len() == c.size {
		tail := c.list.head.prev

		c.list.unlink(tail)
		delete(c.dict, tail.Key)
	}
}

func (c *Cache[K, V]) Len() int {
	return len(c.dict)
}

func (c *Cache[K, V]) Add(key K, value V) {
	if node, ok := c.dict[key]; ok {
		node.Value = value
		c.list.move(node)
		return
	}

	node := &Node[K, V]{Key: key, Value: value}

	c.evict()
	c.list.insert(node)
	c.dict[key] = node
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	if node, ok := c.dict[key]; ok {
		c.list.move(node)
		return node.Value, true
	} else {
		var empty V
		return empty, false
	}
}

func (c *Cache[K, V]) Remove(key K) bool {
	if node, ok := c.dict[key]; ok {
		if node == c.list.head {
			c.list.head = node.next
		}

		c.list.unlink(node)
		delete(c.dict, node.Key)

		return true
	} else {
		return false
	}
}
