package slices_test

import (
	"fmt"
	"testing"

	"github.com/adampresley/adamgokit/slices"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type person struct {
		ID   int
		Name string
	}

	input := []person{
		{ID: 1, Name: "Adam"},
		{ID: 2, Name: "Bob"},
	}

	want := []string{
		`<a href="/person/1">Adam</a>`,
		`<a href="/person/2">Bob</a>`,
	}

	got := slices.Map(input, func(input person, index int) string {
		return fmt.Sprintf(`<a href="/person/%d">%s</a>`, input.ID, input.Name)
	})

	assert.Equal(t, want, got)
}
