package slices_test

import (
	"strconv"
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Run("Convert int to string", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}

		result := slices.Convert(input, func(i int) string {
			return strconv.Itoa(i)
		})

		assert.Equal(t, expected, result)
	})

	t.Run("Convert string to int", func(t *testing.T) {
		input := []string{"10", "20", "30", "40", "50"}
		expected := []int{10, 20, 30, 40, 50}

		result := slices.Convert(input, func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		})

		assert.Equal(t, expected, result)
	})

	t.Run("Convert struct to another struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		type SimplePerson struct {
			FullName string
		}

		input := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35},
		}

		expected := []SimplePerson{
			{FullName: "Alice (30)"},
			{FullName: "Bob (25)"},
			{FullName: "Charlie (35)"},
		}

		result := slices.Convert(input, func(p Person) SimplePerson {
			return SimplePerson{
				FullName: p.Name + " (" + strconv.Itoa(p.Age) + ")",
			}
		})

		assert.Equal(t, expected, result)
	})

	t.Run("Empty slice", func(t *testing.T) {
		var (
			input    []int
			expected = []string{}
		)

		result := slices.Convert(input, func(i int) string {
			return strconv.Itoa(i)
		})

		assert.Equal(t, expected, result)
	})

	t.Run("Nil slice", func(t *testing.T) {
		var expected = []string{}

		result := slices.Convert(nil, func(s *string) string {
			return *s
		})

		assert.Equal(t, expected, result)
	})
}
