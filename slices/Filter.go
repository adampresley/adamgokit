package slices

/*
Filter takes a slice of items, calling filterFn for each item, keeping those that return true.
If filterFn returns false the item is excluded from the result.
*/
func Filter[T comparable](slice []T, filterFn func(item T) bool) []T {
	result := []T{}

	for _, item := range slice {
		if filterFn(item) {
			result = append(result, item)
		}
	}

	return result
}
