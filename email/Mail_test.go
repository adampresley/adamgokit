package email_test

import (
	"testing"

	"github.com/adampresley/adamgokit/email"
	"github.com/stretchr/testify/assert"
)

func TestMail(t *testing.T) {
	from := email.EmailAddress{
		Email: "sender@example.com",
		Name:  "Sender",
	}

	to := []email.EmailAddress{
		{
			Email: "recipient1@example.com",
			Name:  "Recipient 1",
		},
		{
			Email: "recipient2@example.com",
			Name:  "Recipient 2",
		},
	}

	templateData := map[string]any{
		"name":    "John",
		"message": "Hello World",
	}

	mail := email.Mail{
		Body:         "<h1>Test Body</h1>",
		BodyIsHtml:   true,
		From:         from,
		Subject:      "Test Subject",
		Template:     "test-template-id",
		TemplateData: templateData,
		To:           to,
	}

	assert.Equal(t, "<h1>Test Body</h1>", mail.Body)
	assert.True(t, mail.BodyIsHtml)
	assert.Equal(t, from, mail.From)
	assert.Equal(t, "Test Subject", mail.Subject)
	assert.Equal(t, "test-template-id", mail.Template)
	assert.Equal(t, templateData, mail.TemplateData)
	assert.Equal(t, to, mail.To)
	assert.Len(t, mail.To, 2)
}

func TestMailEmpty(t *testing.T) {
	mail := email.Mail{}

	assert.Equal(t, "", mail.Body)
	assert.False(t, mail.BodyIsHtml)
	assert.Equal(t, email.EmailAddress{}, mail.From)
	assert.Equal(t, "", mail.Subject)
	assert.Equal(t, "", mail.Template)
	assert.Nil(t, mail.TemplateData)
	assert.Nil(t, mail.To)
}