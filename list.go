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
	if l.len == 0 {
		l.head = &Node[T]{
			Value: value,
			list:  l,
		}
		l.tail = l.head
	} else {
		oldHead := l.head
		l.head = &Node[T]{
			Value: value,
			next:  oldHead,
			list:  l,
		}
		oldHead.previous = l.head
	}

	l.len++

	return l
}

func (l *List[T]) PushTail(value T) *List[T] {
	if l.len == 0 {
		l.tail = &Node[T]{
			Value: value,
			list:  l,
		}
		l.head = l.tail
	} else {
		oldTail := l.tail
		l.tail = &Node[T]{
			Value:    value,
			previous: oldTail,
			list:     l,
		}
		oldTail.next = l.tail
	}

	l.len++

	return l
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

func (n *Node[T]) RemoveSelf() T {
	if n == n.list.tail {
		return n.list.RemoveTail()
	}

	if n == n.list.head {
		return n.list.RemoveHead()
	}

	n.next.previous = n.previous
	n.previous.next = n.next
	n.list.len--

	return n.Value
}
