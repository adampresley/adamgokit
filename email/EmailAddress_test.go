package email_test

import (
	"testing"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestEmailAddress(t *testing.T) {
	addr := email.EmailAddress{
		Email: "test@example.com",
		Name:  "Test User",
	}

	assert.Equal(t, "test@example.com", addr.Email)
	assert.Equal(t, "Test User", addr.Name)
}

func TestEmailAddressEmpty(t *testing.T) {
	addr := email.EmailAddress{}

	assert.Equal(t, "", addr.Email)
	assert.Equal(t, "", addr.Name)
}