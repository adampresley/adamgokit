package email_test

import (
	"testing"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestIsValidEmailAddress(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.com",
		"user123@example-domain.com",
		"a@b.co",
	}

	for _, emailAddr := range validEmails {
		t.Run("valid_"+emailAddr, func(t *testing.T) {
			assert.True(t, email.IsValidEmailAddress(emailAddr))
		})
	}
}

func TestIsValidEmailAddressInvalid(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"user@",
		"user@@example.com",
		"user..name@example.com",
		"user@.com",
	}

	for _, emailAddr := range invalidEmails {
		t.Run("invalid_"+emailAddr, func(t *testing.T) {
			assert.False(t, email.IsValidEmailAddress(emailAddr))
		})
	}
}