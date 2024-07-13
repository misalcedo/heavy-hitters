package heavy_hitters

// Count is the frequency of an element in a stream along with its estimation error.
type Count struct {
	Count int
	Error int
}
