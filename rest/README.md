# REST Client

This package provides a generic, fluent interface for making HTTP requests and handling responses for RESTful APIs. It simplifies common tasks like setting headers, handling different content types, and managing client configurations.

## Client Configuration

First, create a client configuration. This holds settings that are shared across all requests made with this client, such as the base URL, default headers, and an HTTP client instance.

```go
client := clientoptions.New(
  "https://api.example.com",
  clientoptions.WithDebug(true),
  clientoptions.WithHeaders(map[string]string{
    "Authorization": "Bearer your-token",
  }),
)
```

### Configuration Options

- `WithDebug(bool)`: Enables debug logging, which prints request details (with redacted headers) to the standard logger.
- `WithHeaders(map[string]string)`: Sets headers that will be sent with every request.
- `WithHttpClient(http.Client)`: Allows you to provide a custom `http.Client` instance.

## Making Requests

The client supports `GET`, `POST`, `PUT`, `PATCH`, and `DELETE` methods. These functions are generic and will automatically unmarshal the response body into the type you specify.

### GET and DELETE Requests

Here is an example of making a `GET` request and unmarshalling a JSON response into a struct.

```go
type User struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

// Make the request
user, httpResult, err := rest.Get[User](client, "/users/1")

if err != nil {
  // Handle error
}

fmt.Printf("User: %s", user.Name)
fmt.Printf("Status Code: %d", httpResult.StatusCode)
```

### POST, PUT, and PATCH Requests

To send data, provide an `io.Reader` as the request body.

```go
type NewUser struct {
  Name string `json:"name"`
}

newUser := NewUser{Name: "Adam"}
body, _ := json.Marshal(newUser)

createdUser, _, err := rest.Post[User](client, "/users", bytes.NewReader(body))
```

## Call Options

For individual requests, you can provide call-specific options to override or add to the client configuration.

```go
// Add custom headers and query parameters to a single GET request
user, _, err := rest.Get[User](
  client,
  "/users/1",
  calloptions.WithCallHeaders(map[string]string{
    "X-Request-ID": "12345",
  }),
  calloptions.WithQueryParams(map[string]string{
    "include_details": "true",
  }),
)
```

### Available Call Options

- `WithCallHeaders(map[string]string)`: Adds or overrides headers for a single call.
- `WithQueryParams(map[string]string)`: Appends URL query parameters to the request.
- `WithDebug(bool)`: Overrides the client's debug setting for a single call.

## Response Handling

All request functions return three values:
1.  The unmarshalled result of the type you specified.
2.  An `HttpResult` struct containing the raw response body, status code, and headers.
3.  An `error` if the request failed or if the HTTP status code was not in the 2xx range.

The client automatically handles unmarshalling for `application/json`, `application/xml`, and `text/plain` content types.

## Error Handling

The client will return an error if the HTTP response status code is not a successful 2xx code. Even when an error is returned, the `HttpResult` struct is still populated, allowing you to inspect the response for troubleshooting.

```go
user, httpResult, err := rest.Get[User](client, "/users/999") // Assume this user doesn't exist

if err != nil {
  // Log the error
  slog.Error("API request failed", "error", err)

  // Check the status code and body for more details
  if httpResult.StatusCode == http.StatusNotFound {
    slog.Warn("User not found", "body", string(httpResult.Body))
  }
}
```
