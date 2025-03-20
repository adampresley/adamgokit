package paging_test

import (
	"testing"

	"github.com/adampresley/adamgokit/paging"
	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		Name         string
		Page         int
		TotalItems   int64
		ItemsPerPage int
		Want         paging.Paging
	}{
		{
			Name:         "Expect 3 pages, no previous, has next",
			Page:         1,
			TotalItems:   30,
			ItemsPerPage: 10,
			Want: paging.Paging{
				Page:         1,
				TotalItems:   30,
				ItemsPerPage: 10,
				TotalPages:   3,
				HasNext:      true,
				NextPage:     2,
				HasPrevious:  false,
				PreviousPage: 1,
			},
		},
		{
			Name:         "Expect 3 pages, no next, has previous",
			Page:         3,
			TotalItems:   30,
			ItemsPerPage: 10,
			Want: paging.Paging{
				Page:         3,
				TotalItems:   30,
				ItemsPerPage: 10,
				TotalPages:   3,
				HasNext:      false,
				NextPage:     3,
				HasPrevious:  true,
				PreviousPage: 2,
			},
		},
		{
			Name:         "Expect 1 page, no next, no previous",
			Page:         0,
			TotalItems:   10,
			ItemsPerPage: 10,
			Want: paging.Paging{
				Page:         1,
				TotalItems:   10,
				ItemsPerPage: 10,
				TotalPages:   1,
				HasNext:      false,
				NextPage:     1,
				HasPrevious:  false,
				PreviousPage: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := paging.Calculate(tt.Page, tt.TotalItems, tt.ItemsPerPage)
			assert.Equal(t, tt.Want, got)
		})
	}
}
