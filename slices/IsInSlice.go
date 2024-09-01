package slices

/*
IsInSlice takes an item and a slice of items of type T and
returns true if the item is found in the slice.
*/
func IsInSlice[T comparable](item T, slice []T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}
