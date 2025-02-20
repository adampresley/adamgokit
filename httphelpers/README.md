# HTTP Helpers

This package provides methods of convenience to make reading and writing
HTTP requests and responses easier.

## Requests

### GetFromRequest

**GetFromRequest** attempts to retrieve a value from an HTTP request from all
possible sources, similar to how PHP's `$_REQUEST` array works. Here is the
order of precedence.

1. FORM
2. Query
3. Multipart form data
4. Path

This method uses generics for the type that shold be returned. It supports
the following data types: `int, int32, int64, []int, []int32, []int64, uint, uint32, uint64, []uint, []uint32, []uint64, float32, float64, []float32, []float64, string, []string, bool`.
If a value is not found, the default zero value is returned for that type.

Here is a small sample:

```go
// r is *http.Request
names := httphelpers.GetFromRequest[[]string](r, "names")
age := httphelpers.GetFromRequest[int](r, "age")
```

### GetStringListFromRequest

**GetStringListFromRequest** takes a delimited string from FORM or URL and
returns a split string slice of values.

```go
// Example URL: /?input=1,5,10
inputs := httphelpers.GetStringListFromRequest(r, "input", ",")

// result is []string{"1", "5", "10"}
```

### IsHtmx

**IsHtmx** returns true if the request came from the HTMX library.

```go
// r is a *http.Request struct
isHTMX := httphelpers.IsHtmx(r)
```

### ReadJSONBody

**ReadJSONBody** reads the body content from an http.Request as JSON data into
the provided destination variable.

```go
dest := []string{}

// r is an http.Request
err := httphelpers.ReadJSONBody(r, &dest)
```

## Responses

### WriteJson

**WriteJson** converts any arbitrary structure to JSON and writes it to an HTTP
writer.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

// Here, "w" is an http.ResponseWriter
httphelpers.WriteJson(w, http.StatusOK, output)
```

### JsonOK

**JsonOK** returns a _200 OK_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

httphelpers.JsonOK(w, output)
```

### JsonBadRequest

**JsonBadRequest** returns a _400 Bad Request_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

httphelpers.JsonBadRequest(w, output)
```

### JsonInternalServerError

**JsonBadRequest** returns a _500 Internal Server Error_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

httphelpers.JsonInternalServerError(w, output)
```

### JsonErrorMessage

**JsonErrorMessage** returns a specified status code along with a generic structure containing an error message
in JSON.

```go
httphelpers.JsonErrorMessage(w, http.StatusInternalServerError, "something went wrong")
// The result written is {"message": "something went wrong"}
```

### JsonUnauthorized

**JsonUnauthorized** returns a status code _401 Unauthorized_ along with an arbitrary structure converted to JSON.
in JSON.

```go
httphelpers.JsonErrorMessage(w, http.StatusInternalServerError, "something went wrong")
// The result written is {"message": "something went wrong"}
```
