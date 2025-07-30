package slices

/*
FilterAndMap takes a slice of items, calling f for each item, keeping
transforming those items that return true.
*/
func FilterAndMap[T any, R any](input []T, f func(input T, index int) (R, bool)) []R {
	result := []R{}

	for index, item := range input {
		if newItem, ok := f(item, index); ok {
			result = append(result, newItem)
		}
	}

	return result
}
