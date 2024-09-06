# Slices

This package provides methods for working with slices.

## IsInSlice

A method that returns true if a given item in in a slice.

```go
stringSlice := []string{"A", "B"}
intSlice := []int{1, 5}

// true
stringIsThere := slices.IsInSlice("A", stringSlice)

// false
intIsThere := slices.IsInSlice(4, intSlice)
```
