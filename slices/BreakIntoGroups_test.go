package slices_test

import (
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestBreakIntoGroups(t *testing.T) {
	want := [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	input := []int{1, 2, 3, 4, 5, 6}
	got := slices.BreakIntoGroups(input, 2)

	assert.Equal(t, want, got)
}
