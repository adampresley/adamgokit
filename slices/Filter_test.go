package slices_test

import (
	"strings"
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	input := []string{"Adam", "Bob", "Damn"}
	want := []string{"Adam", "Damn"}

	got := slices.Filter(input, func(item string) bool {
		return strings.Contains(strings.ToLower(item), "dam")
	})

	assert.Equal(t, want, got)
}
