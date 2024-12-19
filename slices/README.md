# Slices

This package provides methods for working with slices.

## BreakIntoGroups

A method to break a slice of items into groups, or an "array of arrays". For example, if you have a slice of strings with 100 items, you can use this method to break that into arrays of 10, resulting in an array of 10 items, each item an array of 10 items.

```go
input := []string{1, 2, 3, 4, 5, 6}
chunks := slices.BreakIntoGroups(input, 2)

// Result:
// [][]string{
//   []string{1, 2},
//   []string{3, 4},
//   []string{5, 6},
// }
```

## Filter

Filter takes a slice of items, calling a filter function for each item, and keeps those that return true. If the filter function returns false, them item is excluded from the resulting slice.

```go
input := []string{"Adam", "Testing", "Bob"}
filtered := slices.Filter(input, func(item string) []string {
  if strings.Contains(item, "dam") {
    return true
  }
  
  return false
})

// filtered == []string{"Adam"}
```

## Find

Find searches a slice for the first item that matches the result of a predicate function. If the function returns true the item is returned. If no item is found then either nil or the type's default value is returned.

```go
input := []int{1, 5, 6, 10}
result := slices.Find(input, func(item int) bool {
  if item == 10 {
    return true
  }

  return false
})
```

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

## Map

Map takes a slice of items (of type T) and function that returns a slice of type R. This is used when you want to transform a slice of items into something else.

```go
type Person struct {
  ID   int
  Name string
  Age  int
}

input := []Person{
  {ID: 1, Name: "Bob", Age: 25},
  {ID: 2, Name: "Tina", Age: 26},
  {ID: 3, Name: "Roger", Age: 35},
}

linksToPeople := slices.Map(input, func(input Person, index int) []string {
  return fmt.Sprintf(`<a href="/person/%d">%s</a>`, input.ID, input.Name)
})

// Result is:
// []string{
//   `<a href="/person/1">Bob</a>`,
//   `<a href="/person/2">Tina</a>`,
//   `<a href="/person/3">Roger</a>`,
// }
```

## Merge

Merge takes two slices, combines their unique values, and returns a new slice.

```go
sliceA := []string{"A", "B", "D"}
sliceB := []string{"B", "E", "F"}

want := []string{"A", "B", "D", "E", "F"}
got := slices.Merge(sliceA, sliceB)

// Result: []string{"A", "B", "D", "E", "F"}
```
