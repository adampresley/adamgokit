package slices_test

import (
	"testing"

	"github.com/adampresley/adamgokit/slices"
)

func TestAllAreEmptyOrHaveValues(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		want   bool
	}{
		{
			name:   "2 empty values returns true",
			inputs: []string{"", ""},
			want:   true,
		},
		{
			name:   "2 filled values returns true",
			inputs: []string{"a", "b"},
			want:   true,
		},
		{
			name:   "many empty values returns true",
			inputs: []string{"", "", "", "", ""},
			want:   true,
		},
		{
			name:   "many filled values returns true",
			inputs: []string{"a", "b", "c", "d", "e"},
			want:   true,
		},
		{
			name:   "1 empty value of two returns false",
			inputs: []string{"", "b"},
			want:   false,
		},
		{
			name:   "several empty values returns false",
			inputs: []string{"a", "", "c", "", "e"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AllAreEmptyOrHaveValues(tt.inputs...)

			if got != tt.want {
				t.Errorf("Test '%s' failed. got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
