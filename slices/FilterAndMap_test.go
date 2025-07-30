package slices

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAndMap(t *testing.T) {
	t.Run("filter and map even integers to strings", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}

		result := FilterAndMap(input, func(item int, index int) (string, bool) {
			if item%2 == 0 {
				return strconv.Itoa(item * 10), true
			}

			return "", false
		})

		expected := []string{"20", "40", "60"}

		assert.Len(t, result, len(expected))
		assert.Equal(t, expected, result)
	})

	t.Run("empty input slice", func(t *testing.T) {
		input := []int{}

		result := FilterAndMap(input, func(item int, index int) (int, bool) {
			return item * 2, true
		})

		assert.Len(t, result, 0)
	})

	t.Run("use index in transformation", func(t *testing.T) {
		input := []string{"apple", "banana", "cherry", "date"}

		result := FilterAndMap(input, func(item string, index int) (string, bool) {
			if index%2 == 0 {
				return item + "_even", true
			}
			return "", false
		})

		expected := []string{"apple_even", "cherry_even"}

		assert.Len(t, result, len(expected))
		assert.Equal(t, expected, result)
	})

	t.Run("nil input", func(t *testing.T) {
		var input []int = nil

		result := FilterAndMap(input, func(item int, index int) (int, bool) {
			return item, true
		})

		expected := []int{}

		assert.Len(t, result, len(expected))
		assert.Equal(t, expected, result)
	})
}

