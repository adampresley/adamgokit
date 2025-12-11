package slices

import "strings"

func AllAreEmptyOrHaveValues(values ...string) bool {
	count := len(values)
	countEmpty := 0
	countValues := 0

	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			countEmpty++
		} else {
			countValues++
		}

	}

	return (countEmpty == count) || (countValues == count)
}
