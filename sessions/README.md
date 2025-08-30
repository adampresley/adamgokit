# Sessions

This package provides components and methods for working with HTTP sessions.

## Postgres Store

This package has a convenience method for creating a session store with
data stored in Postgres. It makes use of the library "github.com/antonlindstrom/pgstore".
Below is a contrived example. Please note: NEVER put secrets in your code!

```go
import "github.com/adampresley/adamgokit/sessions"

dsn := "host=localhost dbname=example user=example password=password port=5432 sslmode=disable"
sessionKey := "password"

store, storeCleanup, err := sessions.NewPGStore(dsn, sessionKey)

if err != nil {
  panic("AHHH!")
}

defer storeCleanup()
```

## Cookie Store

This package provides a convenience method for creating a cookie-based session store
with customizable options. It wraps the Gorilla Sessions cookie store with additional
configuration flexibility.

```go
import "github.com/adampresley/adamgokit/sessions"

sessionKey := "your-secret-key"

// Basic cookie store
store := sessions.NewCookieStore(sessionKey)

// Cookie store with custom options
store := sessions.NewCookieStore(
  sessionKey,
  sessions.WithSecure(true),
  sessions.WithMaxAge(24 * time.Hour),
  sessions.WithSameSite(http.SameSiteStrictMode),
  sessions.WithDomain("example.com"),
  sessions.WithHttpOnly(true),
)
```

Available options:
- `WithSecure(bool)` - Sets whether cookies should only be sent over HTTPS
- `WithMaxAge(time.Duration)` - Sets the maximum age for the cookie
- `WithSameSite(http.SameSite)` - Sets the SameSite attribute for the cookie
- `WithDomain(string)` - Sets the domain for the cookie
- `WithHttpOnly(bool)` - Sets whether the cookie is accessible only through HTTP

## Session Wrapper

This component provides a wrapper around the Gorilla Sessions package. It
provides a type-safe way to get and set session variables.

```go
emailSession := sessions.NewSessionWrapper[string](store, "session", "email")
userIDSession := sessions.NewSessionWrapper[uint](store, "session", "userID")

// r is an *http.Request
email, err := emailSession.Get(r)
userID, err := userIDSession.Get(r)

err = emailSession.Set(r, "new@email.com")
err = userIDSession.Set(r, 25)

err = emailSession.Save(w, r)
err = userIDSession.Save(w, r)
```
