package heavy_hitters

import (
	"github.com/stretchr/testify/require"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestNaiveHeavyHitters(t *testing.T) {
	var hh HeavyHitters[int]

	stream := []int{12, 199997, 30000, 3, 8, 5, 10, 9, 2, 3, 5}
	hh = NewNaive[int]()

	hits := 0

	for _, e := range stream {
		hits++
		hh.Hit(e)
	}

	count := hh.Hit(5)
	require.Equal(t, 3, count.Count)
	require.Equal(t, 0, count.Error)

	top, guaranteed, order := hh.Top(2)
	require.True(t, guaranteed)
	require.True(t, order)
	require.Equal(t, []int{5, 3}, top)

	frequent, guaranteed := hh.Frequent(0.1)
	require.True(t, guaranteed)
	require.Equal(t, []int{5}, frequent)

	require.Equal(t, hits+1, hh.Hits())

	count, found := hh.Get(12)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(199997)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(30000)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(3)
	require.True(t, found)
	require.Equal(t, Count{Count: 2, Error: 0}, count)

	count, found = hh.Get(8)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(5)
	require.True(t, found)
	require.Equal(t, Count{Count: 3, Error: 0}, count)

	count, found = hh.Get(10)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(9)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(2)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = hh.Get(-42)
	require.True(t, found)
	require.Equal(t, Count{Count: 0, Error: 0}, count)
}

func BenchmarkNaive(b *testing.B) {
	seed := time.Now().UTC().UnixNano()

	s := 1.08
	v := 2.0
	imax := uint64(math.MaxUint64)

	generator := rand.NewZipf(rand.New(rand.NewSource(seed)), s, v, imax)
	ss := NewNaive[uint64]()

	b.Run("Hit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ss.Hit(generator.Uint64())
		}
	})

	b.Run("Top", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ss.Top(5)
		}
	})

	top, tGuaranteed, order := ss.Top(5)
	require.Equal(b, []uint64{0, 1, 2, 3, 4}, top)
	require.True(b, tGuaranteed)
	require.True(b, order)

	for _, e := range top {
		_, found := ss.Get(e)
		require.True(b, found)
	}

	b.Run("Frequent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ss.Frequent(0.01)
		}
	})

	frequent, fGuaranteed := ss.Frequent(0.01)
	require.Equal(b, []uint64{0, 1, 2, 3, 4, 5}, frequent)
	require.True(b, fGuaranteed)
}
