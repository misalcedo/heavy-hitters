package heavy_hitters

// List is a doubly-linked list implementation with support for generics.
// Support for generics allows List to avoid the indirection and memory allocations associated with interface types.
type List[T any] struct {
	len  int
	head *Node[T]
	tail *Node[T]
}

// Node is an internal element of a List. Nodes may also be detached from a list.
type Node[T any] struct {
	Value    T
	previous *Node[T]
	next     *Node[T]
	list     *List[T]
}

// Previous node from the point of view of traversing the list from head to tail.
func (n *Node[T]) Previous() *Node[T] {
	return n.previous
}

// Next node from the point of view of traversing the list from head to tail.
func (n *Node[T]) Next() *Node[T] {
	return n.next
}

// NewList creates a new empty List.
func NewList[T any]() *List[T] {
	return &List[T]{}
}

// PushHead appends the given value as the new head of the list.
func (l *List[T]) PushHead(value T) *List[T] {
	node := new(Node[T])
	node.Value = value
	l.PushHeadNode(node)

	return l
}

// PushHeadNode appends the given node as the new head of the list.
// The node will be detached from its previous owning list.
// This enables re-using node pointers to avoid memory allocations and enable control over memory locations.
func (l *List[T]) PushHeadNode(node *Node[T]) {
	if node.list != nil {
		node.RemoveSelf()
	}

	node.list = l

	if l.len == 0 {
		l.head = node
		l.tail = l.head
	} else {
		oldHead := l.head
		l.head = node
		node.next = oldHead
		oldHead.previous = node
	}

	l.len++
}

// PushTail appends the given value as the new tail of the list.
func (l *List[T]) PushTail(value T) *List[T] {
	node := new(Node[T])
	node.Value = value
	l.PushTailNode(node)
	return l
}

// PushTailNode appends the given node as the new tail of the list.
// The node will be detached from its previous owning list.
// This enables re-using node pointers to avoid memory allocations and enable control over memory locations.
func (l *List[T]) PushTailNode(node *Node[T]) {
	if node.list != nil {
		node.RemoveSelf()
	}

	node.list = l

	if l.len == 0 {
		l.tail = node
		l.head = node
	} else {
		oldTail := l.tail
		l.tail = node
		node.previous = oldTail
		oldTail.next = node
	}

	l.len++
}

// Head fetches the head of the list without removing it.
func (l *List[T]) Head() *Node[T] {
	return l.head
}

// Tail fetches the tail of the list without removing it.
func (l *List[T]) Tail() *Node[T] {
	return l.tail
}

// RemoveHead pops the head from the list and returns its value.
func (l *List[T]) RemoveHead() T {
	if l.len == 0 {
		var t T
		return t
	}

	head := l.head

	l.head = l.head.next
	l.len--

	if l.len == 0 {
		l.head = nil
		l.tail = nil
	} else {
		l.head.previous = nil
	}

	return head.Value
}

// RemoveTail pops the tail from the list and returns its value.
func (l *List[T]) RemoveTail() T {
	if l.len == 0 {
		var t T
		return t
	}

	tail := l.tail

	l.tail = l.tail.previous
	l.len--

	if l.len == 0 {
		l.head = nil
		l.tail = nil
	} else {
		l.tail.next = nil
	}

	return tail.Value
}

// Len is the size of the list.
func (l *List[T]) Len() int {
	return l.len
}

// Empty is a predicate that tests if the length of the list is zero.
func (l *List[T]) Empty() bool {
	return l.len == 0
}

// InsertPrevious inserts a node with the given value.
// The new node will be this node's new previous from the point of view of traversing the list from head to tail.
func (n *Node[T]) InsertPrevious(value T) *Node[T] {
	if n == n.list.head {
		n.list.PushHead(value)
		return n.list.head
	}

	node := &Node[T]{
		Value:    value,
		next:     n,
		previous: n.previous,
		list:     n.list,
	}

	n.previous.next = node
	n.previous = node
	n.list.len++

	return node
}

// InsertNext inserts a node with the given value.
// The new node will be this node's new next from the point of view of traversing the list from head to tail.
func (n *Node[T]) InsertNext(value T) *Node[T] {
	if n == n.list.tail {
		n.list.PushTail(value)
		return n.list.tail
	}

	node := &Node[T]{
		Value:    value,
		previous: n,
		next:     n.next,
		list:     n.list,
	}

	n.next.previous = node
	n.next = node
	n.list.len++

	return node
}

// RemoveSelf detaches the current node from its parent list; enabling reuse of the node on other lists (or the same list).
func (n *Node[T]) RemoveSelf() {
	if n == n.list.tail {
		n.list.RemoveTail()
	} else if n == n.list.head {
		n.list.RemoveHead()
	} else {
		if n.next != nil {
			n.next.previous = n.previous
		}

		if n.previous != nil {
			n.previous.next = n.next
		}

		n.list.len--
	}

	n.previous = nil
	n.next = nil
	n.list = nil
}
