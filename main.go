package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

// HeavyHitters provides approximations for finding frequent and top-k elements.
type HeavyHitters[T cmp.Ordered] interface {
	// Hit increments the frequency for the given element, then returns an approximation of the current frequency.
	Hit(T) Count
	// Hits counts the total number of hits for all elements.
	Hits() int
	// Get retrieves the approximated frequency for the given element, with a bounds on the error.
	// The boolean is true the implementation has an approximation for the element's count.
	Get(T) (Count, bool)
	// Frequent finds the set of elements that contribute more than phi * Hits of the total frequency.
	// The slice is returned in descending order of frequency.
	// The boolean is true iff the returned slice is guaranteed to all be frequent elements, irrespective of the errors.
	Frequent(phi float64) ([]T, bool)
	// Top finds the top-k elements seen in the stream.
	// The slice is returned in descending order of frequency.
	// The boolean is true iff the order of the top-k elements is correct and the implementation guarantees they are the actual top-k, irrespective of the errors.
	Top(k int) ([]T, bool)
}

// StreamSummary is a data structure used to implement the [SpaceSaving] algorithm.
// The [SpaceSaving] algorithm reports both frequent and top-k elements with tight guarantees on errors.
//
// [SpaceSaving]: https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
type StreamSummary[T cmp.Ordered] struct {
	hits     int
	elements map[T]*Node[frequencyCounter[T]]
	buckets  *List[frequencyBucket[T]]
}

// frequencyBucket maintains a list of counts with the same frequency.
// The head of the list is the maximum frequency and the tail is the minimum.
type frequencyBucket[T cmp.Ordered] struct {
	count  int
	counts *List[frequencyCounter[T]]
}

type frequencyCounter[T cmp.Ordered] struct {
	Key    T
	Count  int
	Error  int
	bucket *Node[frequencyBucket[T]]
}

type Count struct {
	Count int
	Error int
}

// Hit increments the frequency for the given element, then returns an approximation of the current frequency.
func (s *StreamSummary[T]) Hit(e T) Count {
	s.hits++

	node, monitored := s.elements[e]

	if !monitored {
		// Get the minimum node
		node = s.buckets.Tail().Value.counts.Tail()
		delete(s.elements, node.Value.Key)

		node.Value.Key = e
		node.Value.Error = node.Value.Count
		s.elements[e] = node
	}

	s.incrementCounter(node)

	return Count{Count: node.Value.Count, Error: node.Value.Error}
}

func (s *StreamSummary[T]) incrementCounter(node *Node[frequencyCounter[T]]) {
	oldBucket := node.Value.bucket

	node.Value.bucket = oldBucket.Previous()
	node.Value.Count++

	if node.Value.bucket != nil && node.Value.Count == node.Value.bucket.Value.count {
		node.Value.bucket.Value.counts.PushTailNode(node)
	} else {
		newBucket := oldBucket.InsertPrevious(frequencyBucket[T]{
			count:  node.Value.Count,
			counts: NewList[frequencyCounter[T]](),
		})
		node.Value.bucket = newBucket
		newBucket.Value.counts.PushTailNode(node)
	}

	if oldBucket.Value.counts.Empty() {
		oldBucket.RemoveSelf()
	}
}

// Top finds the top-k elements seen in the stream.
// The slice is returned in descending order of frequency.
// The boolean is true iff the order of the top-k elements is correct and the implementation guarantees they are the actual top-k, irrespective of the errors.
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

// Frequent finds the set of elements that contribute more than phi * Hits of the total frequency.
// The slice is returned in descending order of frequency.
// The boolean is true iff the returned slice is guaranteed to all be frequent elements, irrespective of the errors.
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

// Hits counts the total number of hits for all elements.
func (s *StreamSummary[T]) Hits() int {
	return s.hits
}

// Get retrieves the approximated frequency for the given element, with a bounds on the error.
func (s *StreamSummary[T]) Get(e T) (Count, bool) {
	var count Count

	node, found := s.elements[e]
	if found {
		count = Count{Count: node.Value.Count, Error: node.Value.Error}
	}

	return count, found
}

// NewStreamSummary creates a new instance of a stream summary with the given capacity.
// The error for frequency approximations is guaranteed to be bounded by Hits / capacity.
func NewStreamSummary[T cmp.Ordered](capacity int) *StreamSummary[T] {
	buckets := NewList[frequencyBucket[T]]().PushHead(frequencyBucket[T]{
		counts: NewList[frequencyCounter[T]](),
	})
	bucket := buckets.Tail()

	for i := 0; i < capacity; i++ {
		bucket.Value.counts.PushTail(frequencyCounter[T]{
			bucket: bucket,
		})
	}

	return &StreamSummary[T]{
		elements: make(map[T]*Node[frequencyCounter[T]]),
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
	ss := NewStreamSummary[string](8)
	hits := 0

	for _, e := range s {
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

		fmt.Printf("Element %d is %s with %v\n", i, e, count.Count)
	}

	fmt.Printf("Frequent elements: %v (guaranteed: %v)", frequent, fGuaranteed)
	fmt.Printf("Total hits: %d, Decayed hits: %d", hits, ss.Hits())
}
