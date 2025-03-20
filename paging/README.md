# Paging

A small package to calculate and providing paging information. This package is most useful for communicating page information when querying and working with database records. 

## Calculate
This method calculates several data points around paging information. The returned struct is this.

```go
type Paging struct {
   Page         int
   TotalItems   int64
   ItemsPerPage int
   TotalPages   int
   HasNext      bool
   NextPage     int
   HasPrevious  bool
   PreviousPage int
}
```

To get this data, call the Calculate method like so:

```go
totalItems := 30 // This might come from a database query for example
itemsPerPage := 10
page := 1

pageData := paging.Calculate(page, totalItems, itemsPerPage)

/*
paging = paging.Paging{
   Page: 1,
   TotalItems: 30,
   ItemsPerPage: itemsPerPage,
   TotalPages: 3,
   HasNext: true,
   NextPage: 2,
   HasPrevious: false,
   PreviousPage: 1,
}
*/
```

## Offset
This method returns a offset calculation suitable for use in database queries.

```go
itemsPerPage := 10
page := 2

offset := paging.Offset(page, itemsPerPage)
// offset = 10
```
