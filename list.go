package main

type List[T any] struct {
	len  int
	head *Node[T]
	tail *Node[T]
}

type Node[T any] struct {
	Value    T
	previous *Node[T]
	next     *Node[T]
	list     *List[T]
}

func (n *Node[T]) Previous() *Node[T] {
	return n.previous
}

func (n *Node[T]) Next() *Node[T] {
	return n.next
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) PushHead(value T) *List[T] {
	l.PushHeadNode(&Node[T]{
		Value: value,
	})

	return l
}

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

func (l *List[T]) PushTail(value T) *List[T] {
	l.PushTailNode(&Node[T]{
		Value: value,
	})
	return l
}

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

func (l *List[T]) Head() *Node[T] {
	return l.head
}

func (l *List[T]) Tail() *Node[T] {
	return l.tail
}

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
	}

	return head.Value
}

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
	}

	return tail.Value
}

func (l *List[T]) Len() int {
	return l.len
}

func (l *List[T]) Empty() bool {
	return l.len == 0
}

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
