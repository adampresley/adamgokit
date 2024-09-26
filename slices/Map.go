package slices

/*
Map takes a slice of type []T and a function and returns a slice
of type []R. The function takes two arguments: input of type T,
and an index. The functino should return a type R.
*/
func Map[T any, R any](input []T, f func(input T, index int) R) []R {
	result := []R{}

	for index, item := range input {
		result = append(result, f(item, index))
	}

	return result
}
