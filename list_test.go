package heavy_hitters

import (
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

type Holder struct {
	value int
}

func TestList_Empty(t *testing.T) {
	l := NewList[Holder]()

	require.True(t, l.Empty())
	require.Equal(t, 0, l.Len())
}

func TestList_BasicHead(t *testing.T) {
	var zeroValue *Holder

	l := NewList[*Holder]()

	// Try to break an empty list
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveHead())
	require.Equal(t, 0, l.Len())

	// Try to break a one item list
	l.PushHead(NewHolder(10))
	require.Equal(t, 1, l.Len())
	require.Equal(t, NewHolder(10), l.RemoveHead())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveHead())
	require.Equal(t, 0, l.Len())

	// Mess around
	l.PushHead(NewHolder(10))
	require.Equal(t, 1, l.Len())
	l.PushHead(NewHolder(20))
	require.Equal(t, 2, l.Len())
	l.PushHead(NewHolder(30))
	require.Equal(t, 3, l.Len())
	require.Equal(t, NewHolder(30), l.RemoveHead())
	require.Equal(t, 2, l.Len())
	l.PushHead(NewHolder(40))
	require.Equal(t, 3, l.Len())
	require.Equal(t, NewHolder(40), l.RemoveHead())
	require.Equal(t, 2, l.Len())
	require.Equal(t, NewHolder(20), l.RemoveHead())
	require.Equal(t, 1, l.Len())
	require.Equal(t, NewHolder(10), l.RemoveHead())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveHead())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveHead())
	require.Equal(t, 0, l.Len())
}

func TestList_BasicTail(t *testing.T) {
	var zeroValue *Holder

	l := NewList[*Holder]()

	// Try to break an empty list
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveTail())
	require.Equal(t, 0, l.Len())

	// Try to break a one item list
	l.PushTail(NewHolder(10))
	require.Equal(t, 1, l.Len())
	require.Equal(t, NewHolder(10), l.RemoveTail())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveTail())
	require.Equal(t, 0, l.Len())

	// Mess around
	l.PushTail(NewHolder(10))
	require.Equal(t, 1, l.Len())
	l.PushTail(NewHolder(20))
	require.Equal(t, 2, l.Len())
	l.PushTail(NewHolder(30))
	require.Equal(t, 3, l.Len())
	require.Equal(t, NewHolder(30), l.RemoveTail())
	require.Equal(t, 2, l.Len())
	l.PushTail(NewHolder(40))
	require.Equal(t, 3, l.Len())
	require.Equal(t, NewHolder(40), l.RemoveTail())
	require.Equal(t, 2, l.Len())
	require.Equal(t, NewHolder(20), l.RemoveTail())
	require.Equal(t, 1, l.Len())
	require.Equal(t, NewHolder(10), l.RemoveTail())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveTail())
	require.Equal(t, 0, l.Len())
	require.Equal(t, zeroValue, l.RemoveTail())
	require.Equal(t, 0, l.Len())
}

func TestList_Basic(t *testing.T) {
	var zeroValue *Holder

	m := NewList[*Holder]()

	require.Equal(t, zeroValue, m.RemoveHead())
	require.Equal(t, zeroValue, m.RemoveTail())
	require.Equal(t, zeroValue, m.RemoveHead())
	m.PushHead(NewHolder(1))
	require.Equal(t, NewHolder(1), m.RemoveHead())
	m.PushTail(NewHolder(2))
	m.PushTail(NewHolder(3))
	require.Equal(t, 2, m.Len())
	require.Equal(t, NewHolder(2), m.RemoveHead())
	require.Equal(t, NewHolder(3), m.RemoveHead())
	require.Equal(t, 0, m.Len())
	require.Equal(t, zeroValue, m.RemoveHead())
	m.PushTail(NewHolder(1))
	m.PushTail(NewHolder(3))
	m.PushTail(NewHolder(5))
	m.PushTail(NewHolder(7))
	require.Equal(t, NewHolder(1), m.RemoveHead())
}

func TestList_NonPointer(t *testing.T) {
	n := NewList[Holder]()

	n.PushHead(Holder{2})
	n.PushHead(Holder{3})

	head := n.Head()
	require.Equal(t, Holder{3}, head.Value)
	head.Value.value = 0

	tail := n.Tail()
	require.Equal(t, Holder{2}, tail.Value)
	tail.Value.value = 1

	require.Equal(t, Holder{0}, n.RemoveHead())
	require.Equal(t, Holder{1}, n.RemoveHead())
}

func TestList_Loop(t *testing.T) {
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	actual := make([]int, 0, len(expected))
	l := ListFrom(expected)

	require.Equal(t, len(expected), l.Len())

	for h := l.Head(); h != nil; h = h.Next() {
		actual = append(actual, h.Value)
	}

	require.Equal(t, expected, actual)
}

func TestList_LoopDrain(t *testing.T) {
	var zeroValue int
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	actual := make([]int, 0, len(expected))
	l := ListFrom(expected)

	require.Equal(t, len(expected), l.Len())
	require.NotContains(t, expected, zeroValue)

	for v := l.RemoveHead(); v != zeroValue; v = l.RemoveHead() {
		actual = append(actual, v)
	}

	require.Equal(t, expected, actual)
}

func TestList_LoopTail(t *testing.T) {
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	actual := make([]int, 0, len(expected))
	l := ListFrom(expected)

	require.Equal(t, len(expected), l.Len())

	for t := l.Tail(); t != nil; t = t.Previous() {
		actual = append(actual, t.Value)
	}

	slices.Reverse(expected)
	require.Equal(t, expected, actual)
}

func TestList_LoopTailDrain(t *testing.T) {
	var zeroValue int
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	actual := make([]int, 0, len(expected))
	l := ListFrom(expected)

	require.Equal(t, len(expected), l.Len())
	require.NotContains(t, expected, zeroValue)

	for v := l.RemoveTail(); v != zeroValue; v = l.RemoveTail() {
		actual = append(actual, v)
	}

	slices.Reverse(expected)
	require.Equal(t, expected, actual)
}

func TestList_InsertPreviousHead(t *testing.T) {
	var zeroValue int

	l := ListFrom([]int{2})

	require.Equal(t, 1, l.Len())
	l.Head().InsertPrevious(1)
	require.Equal(t, 2, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.RemoveHead(); v != zeroValue; v = l.RemoveHead() {
		actual = append(actual, v)
	}

	require.Equal(t, []int{1, 2}, actual)
}

func TestList_InsertNextTail(t *testing.T) {
	var zeroValue int

	l := ListFrom([]int{2})

	require.Equal(t, 1, l.Len())
	l.Head().InsertNext(3)
	require.Equal(t, 2, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.RemoveHead(); v != zeroValue; v = l.RemoveHead() {
		actual = append(actual, v)
	}

	require.Equal(t, []int{2, 3}, actual)
}

func TestList_InsertPreviousMiddle(t *testing.T) {
	var zeroValue int

	l := ListFrom([]int{1, 3})

	require.Equal(t, 2, l.Len())
	l.Tail().InsertPrevious(2)
	require.Equal(t, 3, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.RemoveHead(); v != zeroValue; v = l.RemoveHead() {
		actual = append(actual, v)
	}

	require.Equal(t, []int{1, 2, 3}, actual)
}

func TestList_InsertNextMiddle(t *testing.T) {
	var zeroValue int

	l := ListFrom([]int{1, 3})

	require.Equal(t, 2, l.Len())
	l.Head().InsertNext(2)
	require.Equal(t, 3, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.RemoveHead(); v != zeroValue; v = l.RemoveHead() {
		actual = append(actual, v)
	}

	require.Equal(t, []int{1, 2, 3}, actual)
}

func TestList_RemoveSelfHead(t *testing.T) {
	l := ListFrom([]int{1, 3})

	require.Equal(t, 2, l.Len())
	l.Head().RemoveSelf()
	require.Equal(t, 1, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.Tail(); v != nil; v = v.Previous() {
		actual = append(actual, v.Value)
	}

	require.Equal(t, []int{3}, actual)
}

func TestList_RemoveSelfTail(t *testing.T) {
	l := ListFrom([]int{1, 3})

	require.Equal(t, 2, l.Len())
	l.Tail().RemoveSelf()
	require.Equal(t, 1, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.Head(); v != nil; v = v.Next() {
		actual = append(actual, v.Value)
	}

	require.Equal(t, []int{1}, actual)
}

func TestList_RemoveSelf(t *testing.T) {
	l := ListFrom([]int{1, 2, 3})

	require.Equal(t, 3, l.Len())
	l.Head().Next().RemoveSelf()
	require.Equal(t, 2, l.Len())

	actual := make([]int, 0, l.Len())
	for v := l.Tail(); v != nil; v = v.Previous() {
		actual = append(actual, v.Value)
	}

	require.Equal(t, []int{3, 1}, actual)
}

func ListFrom[T any](values []T) *List[T] {
	l := NewList[T]()
	for _, v := range values {
		l.PushTail(v)
	}
	return l
}

func NewHolder(value int) *Holder {
	return &Holder{value}
}
