package email_test

import (
	"testing"
	"time"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestNewMailgunService(t *testing.T) {
	config := &email.Config{
		ApiKey:  "test-api-key",
		Domain:  "example.com",
		Timeout: 30 * time.Second,
	}

	service := email.NewMailgunService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.Config)
	assert.NotNil(t, service.Client)
}