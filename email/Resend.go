package email

import (
	"fmt"

	"github.com/adampresley/adamgokit/slices"
	"github.com/resend/resend-go/v2"
)

type ResendService struct {
	Config *Config
	Client *resend.Client
}

func NewResendService(config *Config) *ResendService {
	return &ResendService{
		Config: config,
		Client: resend.NewClient(config.ApiKey),
	}
}

func (s *ResendService) Send(mail Mail) error {
	msg := &resend.SendEmailRequest{
		From: fmt.Sprintf("%s <%s>", mail.From.Name, mail.From.Email),
		To: slices.Map(mail.To, func(email EmailAddress, index int) string {
			return email.Email
		}),
		Subject: mail.Subject,
		Tags:    []resend.Tag{},
	}

	if mail.BodyIsHtml {
		msg.Html = mail.Body
	} else {
		msg.Text = mail.Body
	}

	_, err := s.Client.Emails.Send(msg)
	return err
}
