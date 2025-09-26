package paging

import "math"

type Paging struct {
	Page         int   `json:"page"`
	TotalItems   int64 `json:"totalItems"`
	ItemsPerPage int   `json:"itemsPerPage"`
	TotalPages   int   `json:"totalPages"`
	HasNext      bool  `json:"hasNext"`
	NextPage     int   `json:"nextPage"`
	HasPrevious  bool  `json:"hasPrevious"`
	PreviousPage int   `json:"previousPage"`
}

/*
Calculate takes page and record total information and returns a struct
that contains data about the number of total pages, and if there are
more pages or previous pages.
*/
func Calculate(page int, totalItems int64, itemsPerPage int) Paging {
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	result := Paging{
		Page:         page,
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		NextPage:     page + 1,
		HasPrevious:  page > 1,
		PreviousPage: page - 1,
	}

	if result.NextPage > totalPages {
		result.NextPage = totalPages
	}

	if result.PreviousPage <= 0 {
		result.PreviousPage = 1
	}

	return result
}

func Offset(page int, itemsPerPage int) int {
	if page <= 0 {
		page = 1
	}

	return (page - 1) * itemsPerPage
}
