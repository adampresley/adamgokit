package auth2

import "fmt"

var (
	ErrSessionKeyNotFound         = fmt.Errorf("session key not found")
	ErrSessionValueNotConvertible = fmt.Errorf("error converting session value to expected type")
)
