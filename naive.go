package heavy_hitters

import (
	"cmp"
	"math"
	"sort"
)

// NaiveHeavyHitters tracks the frequency of elements in a stream by storing all frequencies in a map.
type NaiveHeavyHitters[T cmp.Ordered] struct {
	counts map[T]int
}

func (n NaiveHeavyHitters[T]) Hit(t T) Count {
	count, _ := n.counts[t]
	n.counts[t] = count + 1

	return Count{
		Count: n.counts[t],
	}
}

func (n NaiveHeavyHitters[T]) Hits() int {
	var hits int

	for _, count := range n.counts {
		hits += count
	}

	return hits
}

func (n NaiveHeavyHitters[T]) Get(t T) (Count, bool) {
	count, _ := n.counts[t]
	return Count{
		Count: count,
	}, true
}

func (n NaiveHeavyHitters[T]) Frequent(phi float64) ([]T, bool) {
	threshold := int(math.Ceil(phi * float64(n.Hits())))
	frequent := make([]T, 0)

	for element, count := range n.counts {
		if count > threshold {
			frequent = append(frequent, element)
		}
	}

	sort.Slice(frequent, func(i, j int) bool {
		return n.counts[frequent[i]] > n.counts[frequent[j]]
	})

	return frequent, true
}

func (n NaiveHeavyHitters[T]) Top(k int) ([]T, bool, bool) {
	top := make([]T, 0, len(n.counts))

	for element := range n.counts {
		top = append(top, element)
	}

	sort.Slice(top, func(i, j int) bool {
		return n.counts[top[i]] > n.counts[top[j]]
	})

	return top[0:min(k, len(top))], true, true
}

func NewNaive[T cmp.Ordered]() NaiveHeavyHitters[T] {
	return NaiveHeavyHitters[T]{
		counts: make(map[T]int),
	}
}
