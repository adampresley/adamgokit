package slices_test

import (
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	sliceA := []string{"A", "B", "D"}
	sliceB := []string{"B", "E", "F"}

	want := []string{"A", "B", "D", "E", "F"}
	got := slices.Merge(sliceA, sliceB)

	assert.Equal(t, want, got)
}
