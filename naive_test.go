package heavy_hitters

import (
	"github.com/stretchr/testify/require"
	"testing"
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
