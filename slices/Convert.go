package slices

/*
Converts a slice of any type to a slice of another type using a custom converter function.
*/
func Convert[T any, K any](source []T, converter func(T) K) []K {
	result := make([]K, len(source))

	for index, value := range source {
		result[index] = converter(value)
	}

	return result
}
