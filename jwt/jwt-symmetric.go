package jwt

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type JwtSymmetricService[T any] interface {
	Sign(payload T) (string, error)
	Verify(tokenString string) (T, error)
}

type JwtSymmetric[T any] struct {
	signingKey []byte
}

/*
Creates a new JWT symmetric instance with the given signing key.

Example:

	type MyClaims struct {
		jwt.RegisteredClaims

		UserID    int
		AuthToken string
	}

	signer := NewJwtSymmetric[MyClaims]([]byte("your-secret-key"))
	token := signer.Sign(MyClaims{
		jwt.RegisteredClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "your-issuer",
			Subject: "your-subject",
		},
		UserID: 123,
		AuthToken: "your-auth-token",
	})

	// verify
	token, err := httphelpers.GetAuthorizationBearer(r)

	// handle error

	claims, err := signer.Verify(token)
*/
func NewJwtSymmetric[T any](key []byte) *JwtSymmetric[T] {
	return &JwtSymmetric[T]{
		signingKey: key,
	}
}

func (j *JwtSymmetric[T]) Sign(payload T) (string, error) {
	claims := j.convertToMap(payload)
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(j.signingKey)
}

func (j *JwtSymmetric[T]) Verify(tokenString string) (T, error) {
	var (
		err    error
		token  *gojwt.Token
		result T
	)

	token, err = gojwt.Parse(tokenString, func(t *gojwt.Token) (any, error) {
		return j.signingKey, nil
	})

	if err != nil {
		return result, err
	}

	if claims, ok := token.Claims.(gojwt.MapClaims); ok && token.Valid {
		result, err = j.convertFromMap(claims)

		if err != nil {
			return result, fmt.Errorf("failed to convert token to custom claims struct: %w", err)
		}

		return result, nil
	}

	return result, err
}

func (j *JwtSymmetric[T]) convertToMap(payload T) gojwt.MapClaims {
	claims := make(gojwt.MapClaims)

	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return claims
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Only handle struct types
	if v.Kind() != reflect.Struct {
		return claims
	}

	j.addFieldsToClaims(v, t, claims)
	return claims
}

func (j *JwtSymmetric[T]) addFieldsToClaims(v reflect.Value, t reflect.Type, claims gojwt.MapClaims) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Handle embedded structs (like gojwt.RegisteredClaims)
		if fieldType.Anonymous && field.Kind() == reflect.Struct {
			// Recursively add fields from embedded struct
			j.addFieldsToClaims(field, fieldType.Type, claims)
			continue
		}

		// Use json tag if present, otherwise use field name
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
		}

		claims[fieldName] = field.Interface()
	}
}

func (j *JwtSymmetric[T]) convertFromMap(claims gojwt.MapClaims) (T, error) {
	var result T

	// Get the type of T
	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(&result).Elem()

	// Handle pointer types
	if resultType.Kind() == reflect.Ptr {
		// Create a new instance of the pointed-to type
		elemType := resultType.Elem()
		newElem := reflect.New(elemType)
		resultValue.Set(newElem)
		resultValue = newElem.Elem()
		resultType = elemType
	}

	// Only handle struct types
	if resultType.Kind() != reflect.Struct {
		return result, fmt.Errorf("type T must be a struct, got %s", resultType.Kind())
	}

	for i := 0; i < resultType.NumField(); i++ {
		field := resultValue.Field(i)
		fieldType := resultType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Handle embedded structs (like gojwt.RegisteredClaims)
		if fieldType.Anonymous {
			if err := j.setEmbeddedField(field, claims); err != nil {
				return result, fmt.Errorf("failed to set embedded field %s: %w", fieldType.Name, err)
			}
			continue
		}

		// Determine the field name (use json tag if present)
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
		}

		// Get the value from claims
		claimValue, exists := claims[fieldName]
		if !exists {
			continue
		}

		// Convert and set the field value
		if err := j.setFieldValue(field, claimValue); err != nil {
			return result, fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return result, nil
}

func (j *JwtSymmetric[T]) setEmbeddedField(field reflect.Value, claims gojwt.MapClaims) error {
	fieldType := field.Type()

	// Create a new instance of the embedded struct
	newStruct := reflect.New(fieldType).Elem()

	// Recursively populate the embedded struct fields
	for i := 0; i < fieldType.NumField(); i++ {
		embeddedField := newStruct.Field(i)
		embeddedFieldType := fieldType.Field(i)

		if !embeddedField.CanSet() {
			continue
		}

		// Handle nested embedded structs recursively
		if embeddedFieldType.Anonymous {
			if err := j.setEmbeddedField(embeddedField, claims); err != nil {
				return fmt.Errorf("failed to set nested embedded field %s: %w", embeddedFieldType.Name, err)
			}
			continue
		}

		// Determine the field name (use json tag if present)
		fieldName := embeddedFieldType.Name
		if jsonTag := embeddedFieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				fieldName = jsonTag[:commaIdx]
			} else {
				fieldName = jsonTag
			}
		}

		// Get the value from claims and set it
		if claimValue, exists := claims[fieldName]; exists {
			if err := j.setFieldValue(embeddedField, claimValue); err != nil {
				return fmt.Errorf("failed to set field %s in embedded struct: %w", fieldName, err)
			}
		}
	}

	// Set the populated struct to the field
	field.Set(newStruct)
	return nil
}

func (j *JwtSymmetric[T]) setFieldValue(field reflect.Value, value any) error {
	if value == nil {
		return nil
	}

	valueReflect := reflect.ValueOf(value)
	fieldType := field.Type()

	// Special handling for jwt.NumericDate types
	if fieldType == reflect.TypeOf((*gojwt.NumericDate)(nil)) {
		// Handle *jwt.NumericDate
		if floatVal, ok := value.(float64); ok {
			numericDate := gojwt.NewNumericDate(time.Unix(int64(floatVal), 0))
			field.Set(reflect.ValueOf(numericDate))
			return nil
		}
	} else if fieldType == reflect.TypeOf(gojwt.NumericDate{}) {
		// Handle jwt.NumericDate (non-pointer)
		if floatVal, ok := value.(float64); ok {
			numericDate := *gojwt.NewNumericDate(time.Unix(int64(floatVal), 0))
			field.Set(reflect.ValueOf(numericDate))
			return nil
		}
	}

	// If types match directly, set the value
	if valueReflect.Type().AssignableTo(fieldType) {
		field.Set(valueReflect)
		return nil
	}

	// Try to convert if possible
	if valueReflect.Type().ConvertibleTo(fieldType) {
		field.Set(valueReflect.Convert(fieldType))
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, fieldType)
}
