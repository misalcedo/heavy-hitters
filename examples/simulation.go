package main

import (
	"flag"
	"fmt"
	hh "heavy-hitters"
	"math"
	"math/rand"
	"time"
)

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
	ss := hh.NewStreamSummary[uint64](10)

	for i := 0; i < hits; i++ {
		if zipf {
			ss.Hit(generator.Uint64())
		} else {
			ss.Hit(uint64(rng.NormFloat64()))
		}
	}

	frequent, fGuaranteed := ss.Frequent(0.001)
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
