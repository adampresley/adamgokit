package httphelpers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/stretchr/testify/assert"
)

func TestQueryParamsToString(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "empty query params",
			url:      "/",
			expected: "",
		},
		{
			name:     "single param",
			url:      "/?key=value",
			expected: "key=value",
		},
		{
			name:     "multiple params",
			url:      "/?name=John&age=30&city=Boston",
			expected: "age=30&city=Boston&name=John",
		},
		{
			name:     "param with special characters",
			url:      "/?message=hello%20world&email=test@example.com",
			expected: "email=test%40example.com&message=hello+world",
		},
		{
			name:     "param with multiple values",
			url:      "/?tags=red&tags=blue&tags=green",
			expected: "tags=red%2Cblue%2Cgreen",
		},
		{
			name:     "param with empty value",
			url:      "/?empty=&name=test",
			expected: "empty=&name=test",
		},
		{
			name:     "param with encoded characters",
			url:      "/?data=hello%20world&symbol=%26",
			expected: "data=hello+world&symbol=%26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			result := httphelpers.QueryParamsToString(r)
			assert.Equal(t, tt.expected, result)
		})
	}
}
