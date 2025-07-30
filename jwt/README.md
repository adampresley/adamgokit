# JWT

This is a quick wrapper around the excellent [jwt-go](https://golang-jwt.github.io/jwt/) library. The wrapper just offers a little convience. 

## Example

### Creating a token

```go
import (
	"github.com/adampresley/adamgokit/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	gojwt.RegisteredClaims
	UserID     int
	OtherStuff string
}

jwtService := jwt.NewJwtSymmetric[MyCustomClaims]("my-secret-key")

claims := MyCustomClaims{
	RegisteredClaims: gojwt.RegisteredClaims{
		Issuer:    "issuer",
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		IssuedAt:  gojwt.NewNumericDate(time.Now()),
		NotBefore: gojwt.NewNumericDate(time.Now()),
	},
	UserID:     2,
	OtherStuff: "stuff",
}

if token, err = jwtService.Sign(claims); err != nil {
	// .. handle errors signing
}
```

### Verifying a token

```go
func newApiAuthMiddleware(jwtService jwt.JwtSymmetricService[MyCustomClaims]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				err             error
				tokenFromHeader string
				claims          MyCustomClaims
			)

			if tokenFromHeader, err = httphelpers.GetAuthorizationBearer(r); err != nil {
				// handle error getting auth header
			}

			if claims, err = jwtService.Verify(tokenFromHeader); err != nil {
				// handler error verifying token
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
```
