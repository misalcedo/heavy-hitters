package heavy_hitters

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestSpaceSaving(t *testing.T) {
	var hh HeavyHitters[int]

	stream := []int{12, 199997, 30000, 3, 8, 5, 10, 9, 2, 3, 5}
	hh = NewStreamSummary[int](8)

	hits := 0

	for _, e := range stream {
		hits++
		hh.Hit(e)
	}

	count := hh.Hit(5)
	require.Equal(t, 3, count.Count)
	require.Equal(t, 0, count.Error)

	top, order, guaranteed := hh.Top(2)
	require.True(t, order)
	require.False(t, guaranteed)
	require.Equal(t, []int{5, 2}, top)

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
	require.False(t, found)
	require.Equal(t, Count{Count: 0, Error: 0}, count)

	count, found = hh.Get(2)
	require.True(t, found)
	require.Equal(t, Count{Count: 2, Error: 1}, count)
}

func TestSpaceSaving_ZeroValue(t *testing.T) {
	stream := []int{0, 1, 0}
	hh := NewStreamSummary[int](4)

	for _, e := range stream {
		hh.Hit(e)
	}

	count0, found0 := hh.Get(0)
	require.True(t, found0)
	require.Equal(t, 2, count0.Count)
	require.Equal(t, 0, count0.Error)

	count1, found1 := hh.Get(1)
	require.True(t, found1)
	require.Equal(t, 1, count1.Count)
	require.Equal(t, 0, count1.Error)
}

func BenchmarkSpaceSaving(b *testing.B) {
	seed := time.Now().UTC().UnixNano()

	s := 1.08
	v := 2.0
	imax := uint64(100_000)

	generator := rand.NewZipf(rand.New(rand.NewSource(seed)), s, v, imax)
	ss := NewStreamSummary[uint64](10)

	for i := 0; i < b.N; i++ {
		ss.Hit(generator.Uint64())
	}

	frequent, fGuaranteed := ss.Frequent(0.001)
	require.Equal(b, []uint64{0}, frequent)
	require.True(b, fGuaranteed)

	top, tGuaranteed, order := ss.Top(5)
	require.Equal(b, []uint64{0, 1, 2, 3, 4}, top)
	require.True(b, tGuaranteed)
	require.True(b, order)

	require.Equal(b, b.N, ss.Hits())

	for _, e := range top {
		_, found := ss.Get(e)
		require.True(b, found)
	}
}
