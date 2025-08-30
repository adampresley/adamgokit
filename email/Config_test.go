package email_test

import (
	"testing"
	"time"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := email.Config{
		ApiKey:   "test-api-key",
		Domain:   "example.com",
		Host:     "smtp.example.com",
		Password: "test-password",
		Port:     587,
		Timeout:  30 * time.Second,
		UserName: "test-user",
	}

	assert.Equal(t, "test-api-key", config.ApiKey)
	assert.Equal(t, "example.com", config.Domain)
	assert.Equal(t, "smtp.example.com", config.Host)
	assert.Equal(t, "test-password", config.Password)
	assert.Equal(t, 587, config.Port)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, "test-user", config.UserName)
}