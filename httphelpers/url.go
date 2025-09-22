package httphelpers

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func QueryParamsToString(r *http.Request) string {
	var result strings.Builder
	q := r.URL.Query()

	keys := make([]string, 0, len(q))

	for key := range q {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := q[key]
		result.WriteString(url.QueryEscape(key))
		result.WriteString("=")
		result.WriteString(url.QueryEscape(strings.Join(values, ",")))
		result.WriteString("&")
	}

	return strings.TrimSuffix(result.String(), "&")
}
