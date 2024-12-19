package slices

/*
Merge takes two slices and returns a new, merged slice of unique values.
*/
func Merge[T comparable](sliceA, sliceB []T) []T {
	found := map[T]struct{}{}
	result := []T{}

	for _, a := range sliceA {
		found[a] = struct{}{}
		result = append(result, a)
	}

	for _, b := range sliceB {
		if _, ok := found[b]; !ok {
			result = append(result, b)
		}
	}

	return result
}
