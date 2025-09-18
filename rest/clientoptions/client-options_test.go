package clientoptions_test

import (
	"net/http"
	"testing"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/adampresley/adamgokit/rest/clientoptions"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	headers := map[string]string{"X-Client": "test"}
	mockClient := httphelpers.NewMockHttpClient(t)

	opts := clientoptions.New(
		"http://localhost",
		clientoptions.WithDebug(true),
		clientoptions.WithHeaders(headers),
		clientoptions.WithHttpClient(mockClient),
	)

	assert.Equal(t, "http://localhost", opts.BaseURL)
	assert.True(t, opts.Debug)
	assert.Equal(t, headers, opts.Headers)
	assert.Equal(t, mockClient, opts.HttpClient)
}

func TestNew_Defaults(t *testing.T) {
	opts := clientoptions.New("http://localhost")

	assert.Equal(t, "http://localhost", opts.BaseURL)
	assert.False(t, opts.Debug)
	assert.Nil(t, opts.Headers)
	assert.Equal(t, http.DefaultClient, opts.HttpClient)
}
