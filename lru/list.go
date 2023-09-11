package lru

type Node[K comparable, V any] struct {
	Key        K
	Value      V
	prev, next *Node[K, V]
}

type List[K comparable, V any] struct {
	head *Node[K, V]
}

func (l *List[K, V]) insert(node *Node[K, V]) {
	head := l.head
	l.head = node

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

func (l *List[K, V]) unlink(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
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
		l.unlink(node)
		l.insert(node)
	}
}
