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
