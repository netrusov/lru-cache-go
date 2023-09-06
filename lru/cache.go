package lru

/*
TODO:
	- support concurrent access
	- debug weird behavior when there's only one element in the cache
	- possibly allow `nil` as a cache's head (may fix the bug above), although it will introduce a lot of nil-checks
*/

type Node[K comparable, V any] struct {
	Key        K
	Value      V
	prev, next *Node[K, V]
}

type Cache[K comparable, V any] struct {
	size int
	dict map[K]*Node[K, V]
	head *Node[K, V]
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

func (c *Cache[K, V]) insert(node *Node[K, V]) {
	head := c.head
	c.head = node

	if head == nil {
		node.prev = node
		node.next = node
	} else {
		node.prev = head.prev
		node.next = head

		node.prev.next = node
		node.next.prev = node
	}
}

func (c *Cache[K, V]) unlink(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *Cache[K, V]) move(node *Node[K, V]) {
	if node == c.head {
		return
	}

	if node == c.head.prev {
		// while moving tail to head we can just update the `head` reference
		c.head = node
	} else {
		// otherwise unlink the node and insert it at the beginning
		c.unlink(node)
		c.insert(node)
	}
}

func (c *Cache[K, V]) evict() {
	if c.Len() == c.size {
		tail := c.head.prev

		c.unlink(tail)
		delete(c.dict, tail.Key)
	}
}

func (c *Cache[K, V]) Len() int {
	return len(c.dict)
}

func (c *Cache[K, V]) Add(key K, value V) {
	if node, ok := c.dict[key]; ok {
		node.Value = value
		c.move(node)
		return
	}

	node := &Node[K, V]{Key: key, Value: value}

	c.evict()
	c.insert(node)
	c.dict[key] = node
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	if node, ok := c.dict[key]; ok {
		c.move(node)
		return node.Value, true
	} else {
		var empty V
		return empty, false
	}
}

func (c *Cache[K, V]) Remove(key K) bool {
	if node, ok := c.dict[key]; ok {
		if node == c.head {
			c.head = node.next
		}

		c.unlink(node)
		delete(c.dict, node.Key)

		return true
	} else {
		return false
	}
}
