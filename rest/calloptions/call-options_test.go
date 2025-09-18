package calloptions_test

import (
	"testing"

	"github.com/adampresley/adamgokit/rest/calloptions"
	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	headers := map[string]string{"X-Test": "true"}
	queryParams := map[string]string{"q": "test"}

	opts := &calloptions.CallOptions{}

	calloptions.WithCallHeaders(headers)(opts)
	calloptions.WithDebug(true)(opts)
	calloptions.WithQueryParams(queryParams)(opts)

	assert.True(t, opts.Debug)
	assert.Equal(t, headers, opts.Headers)
	assert.Equal(t, queryParams, opts.QueryParams)
}
