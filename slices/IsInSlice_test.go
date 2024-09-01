package slices_test

import (
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestIsInSlice(t *testing.T) {
	input := []string{"1", "Adam", "Test"}

	got1 := slices.IsInSlice("Adam", input)
	got2 := slices.IsInSlice("Nope", input)

	assert.True(t, got1)
	assert.False(t, got2)
}
