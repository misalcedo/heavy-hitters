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
	capacity int
	elements map[T]Counter[T]
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

	//count, monitored := s.elements[e]

	//if monitored {
	//	s.counts.Remove(count)
	//} else if len(s.elements) >= s.capacity {
	//	count = s.counts.Min()
	//
	//	s.counts.Remove(count)
	//	delete(s.elements, count.Element)
	//	s.elements[e] = count
	//
	//	count.Element = e
	//	count.Error = count.Count
	//} else {
	//	count = &Counter[T]{
	//		Element: e,
	//	}
	//	s.elements[e] = count
	//}
	//
	//count.Count += 1
	//s.counts.Insert(count)
}

func (s *StreamSummary[T]) Top(k int) ([]T, bool) {
	topK := make([]T, 0, k)
	order := true
	guaranteed := false
	minGuaranteedCount := math.MaxInt
	previousGuaranteedCount := math.MaxInt

OuterLoop:
	for b := s.buckets.Tail(); b != nil; b = b.Previous() {
		for c := b.Value.counts.Tail(); c != nil; c = c.Previous() {
			if len(topK) >= k {
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
	for b := s.buckets.Tail(); b != nil; b = b.Previous() {
		for c := b.Value.counts.Tail(); c != nil; c = c.Previous() {
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
	return count, found
}

func New[T cmp.Ordered](capacity int) *StreamSummary[T] {
	return &StreamSummary[T]{
		capacity: capacity,
		elements: make(map[T]Counter[T]),
		buckets:  NewList[Bucket[T]](),
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
