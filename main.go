package main

import (
	"cmp"
	"flag"
	"fmt"
	"math"
	"math/rand"
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
	// The first boolean is true iff the order of the top-k elements is correct and the implementation guarantees they are the actual top-k, irrespective of the errors.
	// The second boolean is true iff the implementation guarantees they are the actual top-k, irrespective of the errors.
	Top(k int) ([]T, bool, bool)
}

// StreamSummary is a data structure used to implement the [SpaceSaving] algorithm.
// The [SpaceSaving] algorithm reports both frequent and top-k elements with tight guarantees on errors.
//
// [SpaceSaving]: https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
type StreamSummary[T cmp.Ordered] struct {
	hits     int
	elements map[T]*Node[frequencyCounter[T]]
	// A list of buckets of counters with the same frequency.
	// The buckets are used to maintain a sorted data structure even in the face of multiple counters with the same frequency.
	// The head of the list is the maximum frequency and the tail is the minimum.
	buckets *List[frequencyBucket[T]]
}

// frequencyBucket maintains a list of counts with the same frequency.
type frequencyBucket[T cmp.Ordered] struct {
	count int
	// A list of counts with the same frequency.
	// The head of the list is the least recently inserted count and the tail is the most recently inserted count.
	counts *List[frequencyCounter[T]]
}

// frequencyCounter counts the frequency of an element in a stream along with its error bounds.
type frequencyCounter[T cmp.Ordered] struct {
	key    T
	count  int
	error  int
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
		// get the node for element with least hits
		// ties can be broken arbitrarily
		node = s.buckets.Tail().Value.counts.Tail()

		// avoid deleting the element from the elements if e is the zero value.
		if node.Value.count > 0 {
			delete(s.elements, node.Value.key)
		}

		// replace the min with e
		node.Value.key = e
		// the error is the value of min
		node.Value.error = node.Value.count
		s.elements[e] = node
	}

	s.incrementCounter(node)

	return Count{Count: node.Value.count, Error: node.Value.error}
}

func (s *StreamSummary[T]) incrementCounter(node *Node[frequencyCounter[T]]) {
	// the current bucket of the node, before incrementing
	oldBucket := node.Value.bucket

	// The previous moves towards the head (assuming head-to-tail traversal).
	// Moving buckets allows us to jump over any other counts with the same frequency.
	node.Value.bucket = oldBucket.Previous()
	node.Value.count += 1

	if node.Value.bucket != nil && node.Value.count == node.Value.bucket.Value.count {
		// If the new bucket exists (the old bucket was not the head), then add this node to the tail.
		// Also, the new bucket's count has to match the count's incremented frequency.
		// Only counts of the same frequency can be in the same bucket.
		node.Value.bucket.Value.counts.PushTailNode(node)
	} else {
		// The old bucket was the head or its count was larger than the node's incremented frequency.
		// Create a new bucket to add this node.
		// The new bucket will either be the head of the list, or be between the old bucket and the old bucket's previous bucket.
		newBucket := oldBucket.InsertPrevious(frequencyBucket[T]{
			count:  node.Value.count,
			counts: NewList[frequencyCounter[T]](),
		})
		// Add this node to the new bucket.
		node.Value.bucket = newBucket
		newBucket.Value.counts.PushTailNode(node)
	}

	// If the old bucket is empty, remove it from the list of buckets.
	if oldBucket.Value.counts.Empty() {
		oldBucket.RemoveSelf()
	}
}

// Top finds the top-k elements seen in the stream.
// The slice is returned in descending order of frequency.
// The first boolean is true iff the order of the top-k elements is correct and the implementation guarantees they are the actual top-k, irrespective of the errors.
// The second boolean is true iff the implementation guarantees they are the actual top-k, irrespective of the errors.
func (s *StreamSummary[T]) Top(k int) ([]T, bool, bool) {
	topK := make([]T, 0, k)
	order := true
	guaranteed := false
	minGuaranteedCount := math.MaxInt
	previousGuaranteedCount := math.MaxInt

OuterLoop:
	for b := s.buckets.Head(); b != nil; b = b.Next() {
		if b.Value.count == 0 {
			continue
		}

		for c := b.Value.counts.Head(); c != nil; c = c.Next() {
			if len(topK) >= k {
				guaranteed = c.Value.count <= minGuaranteedCount
				break OuterLoop
			}

			topK = append(topK, c.Value.key)
			guaranteedCount := c.Value.count - c.Value.error
			minGuaranteedCount = min(minGuaranteedCount, guaranteedCount)
			order = order && (guaranteedCount <= previousGuaranteedCount)

			previousGuaranteedCount = guaranteedCount
		}
	}

	return topK, order, guaranteed
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
		if b.Value.count <= threshold {
			// all counts in the same bucket have the same frequency, so we only need to test this predicate once per bucket.
			break OuterLoop
		}

		if b.Value.count == 0 {
			continue
		}

		for c := b.Value.counts.Head(); c != nil; c = c.Next() {
			frequent = append(frequent, c.Value.key)
			guaranteed = guaranteed && ((c.Value.count - c.Value.error) >= threshold)
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
		count = Count{Count: node.Value.count, Error: node.Value.error}
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
	var hits int
	var seed int64
	var zipf bool

	// See https://en.wikipedia.org/wiki/Zipf%27s_law
	var a, b float64
	var imax uint64

	flag.IntVar(&hits, "h", 1_000_000, "number of hits in total")
	flag.BoolVar(&zipf, "z", true, "use a zipf distribution")
	flag.Int64Var(&seed, "seed", time.Now().UTC().UnixNano(), "seed for the random number generator")
	flag.Float64Var(&a, "a", 1.08, "a parameter for Zipf's law")
	flag.Float64Var(&b, "b", 2, "b parameter for Zipf's law'")
	flag.Uint64Var(&imax, "imax", math.MaxUint64, "imax parameter for Zipf generator")
	flag.Parse()

	fmt.Printf("Running simulation with seed %d.\n", seed)

	rng := rand.New(rand.NewSource(seed))
	generator := rand.NewZipf(rng, a, b, imax)
	start := time.Now()
	ss := NewStreamSummary[uint64](10)

	for i := 0; i < hits; i++ {
		if zipf {
			ss.Hit(generator.Uint64())
		} else {
			ss.Hit(uint64(rng.NormFloat64()))
		}
	}

	frequent, fGuaranteed := ss.Frequent(0.1)
	top, tGuaranteed, order := ss.Top(5)

	fmt.Printf("Elapsed: %s\n", time.Since(start))
	fmt.Printf("Total hits: %d, summarized hits: %d\n", hits, ss.Hits())
	fmt.Printf("Frequent elements: %v (guaranteed: %v)\n", frequent, fGuaranteed)
	fmt.Printf("Top elements (guaranteed: %v, order: %v): %v\n", tGuaranteed, order, top)

	for i, e := range top {
		count, found := ss.Get(e)

		if !found {
			panic("unable to find element")
		}

		fmt.Printf("Top-%d is %d: {count: %d, error: %d}}\n", i+1, e, count.Count, count.Error)
	}
}
