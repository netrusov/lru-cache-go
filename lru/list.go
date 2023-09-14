package lru

type Node[K comparable, V any] struct {
	Key        K
	Value      V
	prev, next *Node[K, V]
}

type List[K comparable, V any] struct {
	head *Node[K, V]
	len  int
}

func (l *List[K, V]) insert(node *Node[K, V]) {
	if l.head == nil {
		node.prev = node
		node.next = node
	} else {
		node.prev = l.head.prev
		node.next = l.head

		node.prev.next = node
		node.next.prev = node
	}

	l.head = node

	l.len++
}

func (l *List[K, V]) remove(node *Node[K, V]) {
	if l.len == 1 {
		l.head = nil
	} else {
		if node == l.head {
			l.head = node.next
		}

		node.prev.next = node.next
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = nil

	l.len--
}

func (l *List[K, V]) move(node *Node[K, V]) {
	if node == l.head {
		return
	}

	if node == l.head.prev {
		// while moving tail to head we can just update the `head` reference
		l.head = node
	} else {
		// otherwise unlink the node and insert it at the beginning
		l.remove(node)
		l.insert(node)
	}
}
