package slices

/*
Converts a slice of any type to a slice of another type using a custom converter function.
If source is nil or empty, the function returns an empty slice of the target type.
*/
func Convert[T any, K any](source []T, converter func(T) K) []K {
	result := make([]K, len(source))

	for index, value := range source {
		result[index] = converter(value)
	}

	return result
}
