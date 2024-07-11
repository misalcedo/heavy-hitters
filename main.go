package main

import (
	"fmt"
	goset "github.com/hashicorp/go-set/v2"
	"math"
	"os"
	"strings"
	"time"
)

// StreamSummary implements the [SpaceSaving] algorithm.
//
// [SpaceSaving]: https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
type StreamSummary[T goset.GoType] struct {
	hits     int
	capacity int
	elements map[T]*Counter[T]
	counts   *goset.TreeSet[*Counter[T]]
}

type Counter[T goset.GoType] struct {
	Element T
	Count   int
	Error   int
}

func Compare[T goset.GoType](a, b *Counter[T]) int {
	if a.Count != b.Count {
		return goset.Compare(a.Count, b.Count)
	}

	if a.Error != b.Error {
		return goset.Compare(a.Error, b.Error)
	}

	return goset.Compare(a.Element, b.Element)
}

func (s *StreamSummary[T]) Hit(e T) {
	s.hits++

	count, monitored := s.elements[e]

	if monitored {
		s.counts.Remove(count)
	} else if len(s.elements) >= s.capacity {
		count = s.counts.Min()

		s.counts.Remove(count)
		delete(s.elements, count.Element)
		s.elements[e] = count

		count.Element = e
		count.Error = count.Count
	} else {
		count = &Counter[T]{
			Element: e,
		}
		s.elements[e] = count
	}

	count.Count += 1
	s.counts.Insert(count)
}

func (s *StreamSummary[T]) Top(k int) ([]T, bool) {
	topK := make([]T, 0, k)
	order := true
	guaranteed := false
	minGuaranteedCount := math.MaxInt

	requestedLen := k + 1
	topCounters := s.counts.TopK(requestedLen)
	actualLen := len(topCounters)

	for i := actualLen - 1; i >= 0; i-- {
		c := topCounters[i]

		guaranteedCount := c.Count - c.Error
		minGuaranteedCount = min(minGuaranteedCount, guaranteedCount)

		if len(topK) < k {
			topK = append(topK, c.Element)

			if i > 0 {
				order = order && (guaranteedCount >= topCounters[i-1].Count)
			}
		}
	}

	if actualLen == requestedLen {
		guaranteed = topCounters[0].Count <= minGuaranteedCount
	}

	return topK, guaranteed && order
}

func (s *StreamSummary[T]) Frequent(phi float64) ([]T, bool) {
	threshold := int(math.Ceil(phi * float64(s.hits)))
	frequent := make([]T, 0)
	guaranteed := true

	s.counts.ForEach(func(c *Counter[T]) bool {
		frequent = append(frequent, c.Element)
		guaranteed = guaranteed && ((c.Count - c.Error) >= threshold)

		return c.Count <= threshold
	})

	return frequent, guaranteed
}

func (s *StreamSummary[T]) Hits() int {
	return s.hits
}

func (s *StreamSummary[T]) Get(e T) (Counter[T], bool) {
	result := Counter[T]{
		Element: e,
	}

	count, found := s.elements[e]
	if found {
		result.Count = count.Count
		result.Error = count.Error
	}

	return result, found
}

func New[T goset.GoType](capacity int) *StreamSummary[T] {
	return &StreamSummary[T]{
		capacity: capacity,
		elements: make(map[T]*Counter[T]),
		counts:   goset.NewTreeSet[*Counter[T]](Compare[T]),
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

	fmt.Printf("Frequent elements: %v (guaranteed: %s)", frequent, fGuaranteed)
	fmt.Printf("Total hits: %d, Decayed hits: %d", hits, ss.Hits())
}
