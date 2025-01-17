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

/*
FindWithIndex searches for an item in slice []T by calling a predicate function.
If the function returns true then the item and its index is returned.
*/
func FindWithIndex[T any](slice []T, f func(item T) bool) (T, int) {
	var result T

	for idx, i := range slice {
		if f(i) {
			return i, idx
		}
	}

	return result, -1
}
