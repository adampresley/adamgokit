package httphelpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetFromRequest[T RequestTypes](r *http.Request, name string) T {
	var (
		err  error
		cast T
	)

	switch any(cast).(type) {
	default:
		result := r.FormValue(name)

		if result == "" {
			result = r.PathValue(name)
		}

		return any(result).(T)

	case []string:
		result := r.Form[name]
		return any(result).(T)

	case bool:
		temp := r.FormValue(name)
		result := false

		if temp == "" {
			temp = r.PathValue(name)
		}

		if result, err = strconv.ParseBool(temp); err != nil {
			result = false
		}

		return any(result).(T)

	case int:
		result := getInt(r, name, 64)
		return any(int(result)).(T)

	case int32:
		result := getInt(r, name, 32)
		return any(int32(result)).(T)

	case int64:
		result := getInt(r, name, 64)
		return any(result).(T)

	case uint:
		result := getUint(r, name, 64)
		return any(uint(result)).(T)

	case uint32:
		result := getUint(r, name, 32)
		return any(uint32(result)).(T)

	case uint64:
		result := getUint(r, name, 64)
		return any(result).(T)

	case []int:
		temp := getIntSlice(r, name, 64)
		result := []int{}

		for _, v := range temp {
			result = append(result, int(v))
		}

		return any(result).(T)

	case []int32:
		temp := getIntSlice(r, name, 32)
		result := []int32{}

		for _, v := range temp {
			result = append(result, int32(v))
		}

		return any(result).(T)

	case []int64:
		result := getIntSlice(r, name, 64)
		return any(result).(T)

	case []uint:
		temp := getUintSlice(r, name, 64)
		result := []uint{}

		for _, v := range temp {
			result = append(result, uint(v))
		}

		return any(result).(T)

	case []uint32:
		temp := getUintSlice(r, name, 32)
		result := []uint32{}

		for _, v := range temp {
			result = append(result, uint32(v))
		}

		return any(result).(T)

	case []uint64:
		result := getUintSlice(r, name, 64)
		return any(result).(T)

	case float32:
		result := getFloat(r, name, 32)
		return any(float32(result)).(T)

	case []float32:
		temp := getFloatSlice(r, name, 32)
		result := []float32{}

		for _, v := range temp {
			result = append(result, float32(v))
		}

		return any(result).(T)

	case float64:
		result := getFloat(r, name, 64)
		return any(result).(T)

	case []float64:
		result := getFloatSlice(r, name, 64)
		return any(result).(T)
	}
}

/*
GetAuthorizationBearer returns the token portion of a Bearer authorization header.
If the header is missing or malformed, an error is returned.
*/
func GetAuthorizationBearer(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	bearerPart := strings.TrimSpace(authHeader)
	bearerParts := strings.Split(bearerPart, " ")

	if len(bearerParts) < 2 {
		return "", fmt.Errorf("invalid bearer authorization header")
	}

	return strings.TrimSpace(bearerParts[1]), nil
}

/*
GetStringListFromRequest returns a string slice from a delimited string on
form or query param.
*/
func GetStringListFromRequest(r *http.Request, name, seperator string) []string {
	var (
		value []string
	)

	values := r.FormValue(name)
	value = strings.Split(values, seperator)
	return value
}

/*
IsHtmx returns true if the request came from the Htmx library.
*/
func IsHtmx(r *http.Request) bool {
	return r.Header.Get("Hx-Request") != ""
}

/*
ReadJSONBody reads the body content from an http.Request as JSON data into
dest.
*/
func ReadJSONBody(r *http.Request, dest interface{}) error {
	var (
		err error
		b   []byte
	)

	if b, err = io.ReadAll(r.Body); err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	if err = json.Unmarshal(b, &dest); err != nil {
		return fmt.Errorf("error unmarshaling body to destination: %w", err)
	}

	return nil
}

func getInt(r *http.Request, name string, size int) int64 {
	var (
		err    error
		result int64
	)

	valueString := r.FormValue(name)

	if result, err = strconv.ParseInt(valueString, 10, size); err != nil {
		valueString = r.PathValue(name)

		if result, err = strconv.ParseInt(valueString, 10, size); err != nil {
			result = 0
		}
	}

	return result
}

func getIntSlice(r *http.Request, name string, size int) []int64 {
	var (
		err    error
		values []string
		result []int64
		temp   int64
	)

	values = r.Form[name]

	if len(values) > 0 {
		for _, v := range values {
			if temp, err = strconv.ParseInt(v, 10, size); err == nil {
				result = append(result, temp)
			}
		}
	}

	return result
}

func getUint(r *http.Request, name string, size int) uint64 {
	var (
		err    error
		result uint64
	)

	valueString := r.FormValue(name)

	if result, err = strconv.ParseUint(valueString, 10, size); err != nil {
		valueString = r.PathValue(name)

		if result, err = strconv.ParseUint(valueString, 10, size); err != nil {
			result = 0
		}
	}

	return result
}

func getUintSlice(r *http.Request, name string, size int) []uint64 {
	var (
		err    error
		values []string
		result []uint64
		temp   uint64
	)

	values = r.Form[name]

	if len(values) > 0 {
		for _, v := range values {
			if temp, err = strconv.ParseUint(v, 10, size); err == nil {
				result = append(result, temp)
			}
		}
	}

	return result
}

func getFloat(r *http.Request, name string, size int) float64 {
	var (
		err    error
		result float64
	)

	valueString := r.FormValue(name)

	if result, err = strconv.ParseFloat(valueString, size); err != nil {
		valueString = r.PathValue(name)

		if result, err = strconv.ParseFloat(valueString, size); err != nil {
			result = 0.0
		}
	}

	return result
}

func getFloatSlice(r *http.Request, name string, size int) []float64 {
	var (
		err    error
		values []string
		result []float64
		temp   float64
	)

	values = r.Form[name]

	if len(values) > 0 {
		for _, v := range values {
			if temp, err = strconv.ParseFloat(v, size); err == nil {
				result = append(result, temp)
			}
		}
	}

	return result
}
