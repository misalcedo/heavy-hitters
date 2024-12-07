package misra_gries

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMisraGries(t *testing.T) {
	hh, err := NewMisraGries(0.125)
	require.NoError(t, err)

	stream := []string{"12", "199997", "30000", "3", "8", "5", "10", "9", "2", "3", "5"}
	hits := 0

	for _, e := range stream {
		hits++
		hh.Hit(e)
	}

	require.Equal(t, hits, hh.Hits())

	low, high := hh.Query("0")
	require.Equal(t, 0, low)
	require.Equal(t, 2, high)

	low, high = hh.Query("3")
	require.Equal(t, 1, low)
	require.Equal(t, 3, high)
}
