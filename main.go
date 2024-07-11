package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

// StreamSummary implements the [SpaceSaving] algorithm.
//
// [SpaceSaving]: https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
type StreamSummary[T cmp.Ordered] struct {
	hits     int
	elements map[T]*Node[Counter[T]]
	buckets  *List[Bucket[T]]
}

type Bucket[T cmp.Ordered] struct {
	count  int
	counts *List[Counter[T]]
}

type Counter[T cmp.Ordered] struct {
	Key    T
	Count  int
	Error  int
	bucket *Node[Bucket[T]]
}

func (s *StreamSummary[T]) Hit(e T) {
	s.hits++

	node, monitored := s.elements[e]

	if !monitored {
		// Get the minimum node
		node = s.buckets.Tail().Value.counts.Tail()
		count := node.Value
		delete(s.elements, count.Key)

		count.Key = e
		count.Error = count.Count
		s.elements[e] = node
	}

	s.incrementCounter(node)
}

func (s *StreamSummary[T]) incrementCounter(node *Node[Counter[T]]) {
	count := node.Value
	bucket := count.bucket
	count.bucket = bucket.Next()

	node.RemoveSelf()
}

func (s *StreamSummary[T]) Top(k int) ([]T, bool) {
	topK := make([]T, 0, k)
	order := true
	guaranteed := false
	minGuaranteedCount := math.MaxInt
	previousGuaranteedCount := math.MaxInt

OuterLoop:
	for b := s.buckets.Head(); b != nil; b = b.Next() {
		for c := b.Value.counts.Head(); c != nil; c = c.Next() {
			if len(topK) >= k {
				guaranteed = c.Value.Count <= minGuaranteedCount
				break OuterLoop
			}

			topK = append(topK, c.Value.Key)
			guaranteedCount := c.Value.Count - c.Value.Error
			minGuaranteedCount = min(minGuaranteedCount, guaranteedCount)
			order = order && (guaranteedCount <= previousGuaranteedCount)

			previousGuaranteedCount = guaranteedCount
		}
	}

	return topK, guaranteed && order
}

func (s *StreamSummary[T]) Frequent(phi float64) ([]T, bool) {
	threshold := int(math.Ceil(phi * float64(s.hits)))
	frequent := make([]T, 0)
	guaranteed := true

OuterLoop:
	for b := s.buckets.Head(); b != nil; b = b.Next() {
		for c := b.Value.counts.Head(); c != nil; c = c.Next() {
			if b.Value.count <= threshold {
				break OuterLoop
			}

			frequent = append(frequent, c.Value.Key)
			guaranteed = guaranteed && ((c.Value.Count - c.Value.Error) >= threshold)
		}
	}

	return frequent, guaranteed
}

func (s *StreamSummary[T]) Hits() int {
	return s.hits
}

func (s *StreamSummary[T]) Get(e T) (Counter[T], bool) {
	count, found := s.elements[e]
	return count.Value, found
}

func New[T cmp.Ordered](capacity int) *StreamSummary[T] {
	buckets := NewList[Bucket[T]]().PushHead(Bucket[T]{
		counts: NewList[Counter[T]](),
	})
	bucket := buckets.Tail()

	for i := 0; i < capacity; i++ {
		bucket.Value.counts.PushTail(Counter[T]{
			bucket: bucket,
		})
	}

	return &StreamSummary[T]{
		elements: make(map[T]*Node[Counter[T]]),
		buckets:  buckets,
	}
}

func main() {
	path := os.Args[1]
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	start := time.Now()

	s := strings.Fields(string(contents))
	ss := New[string](8)
	hits := 0

	for _, e := range s {
		if hits > 20 {
			break
		}

		hits++
		ss.Hit(e)
	}

	top, tGuaranteed := ss.Top(2)
	if !tGuaranteed {
		panic("unable to guarantee top hitters")
	}

	frequent, fGuaranteed := ss.Frequent(0.1)

	fmt.Printf("Elapsed: %s\n", time.Since(start))
	fmt.Printf("Top elements: %v\n", top)

	for i, e := range top {
		count, found := ss.Get(e)
		if !found {
			panic("unable to find element")
		}

		fmt.Printf("Element %d is %s with %v", i, e, count)
	}

	fmt.Printf("Frequent elements: %v (guaranteed: %v)", frequent, fGuaranteed)
	fmt.Printf("Total hits: %d, Decayed hits: %d", hits, ss.Hits())
}
