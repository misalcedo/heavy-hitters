package heavy_hitters

import (
	"cmp"
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
