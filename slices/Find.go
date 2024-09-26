package slices

/*
Find searches for an item in slice []T by calling a predicate function.
If the function returns true then the item is returned.
*/
func Find[T any](slice []T, f func(item T) bool) T {
	var result T

	for _, i := range slice {
		if f(i) {
			return i
		}
	}

	return result
}
