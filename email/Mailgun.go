package email

import (
	"context"

	"github.com/adampresley/adamgokit/slices"
	"github.com/mailgun/mailgun-go/v5"
)

type MailgunService struct {
	Config *Config
	Client *mailgun.Client
}

func NewMailgunService(config *Config) *MailgunService {
	return &MailgunService{
		Config: config,
		Client: mailgun.NewMailgun(config.ApiKey),
	}
}

func (s *MailgunService) Send(mail Mail) error {
	var message mailgun.Message
	body := ""

	if !mail.BodyIsHtml {
		body = mail.Body
	}

	if mail.BodyIsHtml {
		message = mailgun.NewMessage(
			s.Config.Domain,
			mail.From.Email,
			mail.Subject,
			"",
			slices.Map(mail.To, func(input EmailAddress, index int) string {
				return input.Email
			})...,
		)

		message.SetHTML(body)

	} else {
		message = mailgun.NewMessage(
			s.Config.Domain,
			mail.From.Email,
			mail.Subject,
			body,
			slices.Map(mail.To, func(input EmailAddress, index int) string {
				return input.Email
			})...,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.Config.Timeout)
	defer cancel()

	_, err := s.Client.Send(ctx, message)
	return err
}
