package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSpaceSaving(t *testing.T) {
	stream := []int{12, 199997, 30000, 3, 8, 5, 10, 9, 2, 3, 5}
	ss := NewStreamSummary[int](8)
	hits := 0

	for _, e := range stream {
		hits++
		ss.Hit(e)
	}

	count := ss.Hit(5)
	require.Equal(t, 3, count.Count)
	require.Equal(t, 0, count.Error)

	top, guaranteed := ss.Top(2)
	require.False(t, guaranteed)
	require.Equal(t, []int{5, 2}, top)

	frequent, guaranteed := ss.Frequent(0.1)
	require.True(t, guaranteed)
	require.Equal(t, []int{5}, frequent)

	require.Equal(t, hits+1, ss.Hits())

	count, found := ss.Get(12)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = ss.Get(199997)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = ss.Get(30000)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = ss.Get(3)
	require.True(t, found)
	require.Equal(t, Count{Count: 2, Error: 0}, count)

	count, found = ss.Get(8)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = ss.Get(5)
	require.True(t, found)
	require.Equal(t, Count{Count: 3, Error: 0}, count)

	count, found = ss.Get(10)
	require.True(t, found)
	require.Equal(t, Count{Count: 1, Error: 0}, count)

	count, found = ss.Get(9)
	require.False(t, found)
	require.Equal(t, Count{Count: 0, Error: 0}, count)

	count, found = ss.Get(2)
	require.True(t, found)
	require.Equal(t, Count{Count: 2, Error: 1}, count)
}
