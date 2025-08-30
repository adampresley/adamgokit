package email_test

import (
	"testing"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestNewSendGridService(t *testing.T) {
	config := &email.Config{
		ApiKey: "test-api-key",
		Domain: "example.com",
	}

	service := email.NewSendGridService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.Config)
}

func TestSendGridServiceSendMissingFrom(t *testing.T) {
	config := &email.Config{
		ApiKey: "test-api-key",
		Domain: "example.com",
	}

	service := email.NewSendGridService(config)

	mail := email.Mail{
		Subject: "Test Subject",
		To: []email.EmailAddress{
			{Email: "recipient@example.com", Name: "Recipient"},
		},
		Template: "test-template",
	}

	err := service.Send(mail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "from email address is required")
}

func TestSendGridServiceSendMissingTo(t *testing.T) {
	config := &email.Config{
		ApiKey: "test-api-key",
		Domain: "example.com",
	}

	service := email.NewSendGridService(config)

	mail := email.Mail{
		From: email.EmailAddress{
			Email: "sender@example.com",
			Name:  "Sender",
		},
		Subject:  "Test Subject",
		Template: "test-template",
	}

	err := service.Send(mail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "to email address is required")
}

func TestSendGridServiceSendEmptyTo(t *testing.T) {
	config := &email.Config{
		ApiKey: "test-api-key",
		Domain: "example.com",
	}

	service := email.NewSendGridService(config)

	mail := email.Mail{
		From: email.EmailAddress{
			Email: "sender@example.com",
			Name:  "Sender",
		},
		Subject:  "Test Subject",
		To:       []email.EmailAddress{},
		Template: "test-template",
	}

	err := service.Send(mail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "to email address is required")
}