package email_test

import (
	"testing"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestNewSmtpMailService(t *testing.T) {
	config := &email.Config{
		Host:     "smtp.example.com",
		Port:     587,
		UserName: "test-user",
		Password: "test-password",
	}

	service := email.NewSmtpMailService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.Config)
	assert.NotNil(t, service.Dialer)
	assert.Equal(t, "smtp.example.com", service.Dialer.Host)
	assert.Equal(t, 587, service.Dialer.Port)
	assert.Equal(t, "test-user", service.Dialer.Username)
	assert.Equal(t, "test-password", service.Dialer.Password)
}

func TestNewSmtpMailServiceNoCredentials(t *testing.T) {
	config := &email.Config{
		Host: "smtp.example.com",
		Port: 587,
	}

	service := email.NewSmtpMailService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.Config)
	assert.NotNil(t, service.Dialer)
	assert.Equal(t, "smtp.example.com", service.Dialer.Host)
	assert.Equal(t, 587, service.Dialer.Port)
	assert.Equal(t, "", service.Dialer.Username)
	assert.Equal(t, "", service.Dialer.Password)
}

func TestNewSmtpMailServicePartialCredentials(t *testing.T) {
	config := &email.Config{
		Host:     "smtp.example.com",
		Port:     587,
		UserName: "test-user",
	}

	service := email.NewSmtpMailService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.Config)
	assert.NotNil(t, service.Dialer)
	assert.Equal(t, "smtp.example.com", service.Dialer.Host)
	assert.Equal(t, 587, service.Dialer.Port)
	assert.Equal(t, "test-user", service.Dialer.Username)
	assert.Equal(t, "", service.Dialer.Password)
}