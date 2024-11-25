package slices_test

import (
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	input := []int{1, 5, 12, 10, 2, 11}
	want := 12

	got := slices.Find(input, func(item int) bool {
		return item > 10
	})

	assert.Equal(t, want, got)
}
