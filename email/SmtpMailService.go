package email

import (
	"gopkg.in/gomail.v2"
)

/*
SmtpMailService provides methods for working with email
*/
type SmtpMailService struct {
	Config *Config
	Dialer *gomail.Dialer
}

/*
NewSmtpMailService creates a new instance of MailService
*/
func NewSmtpMailService(config *Config) *SmtpMailService {
	result := &SmtpMailService{
		Config: config,
		Dialer: &gomail.Dialer{
			Host: config.Host,
			Port: config.Port,
		},
	}

	if config.UserName != "" {
		result.Dialer.Username = config.UserName
	}

	if config.Password != "" {
		result.Dialer.Password = config.Password
	}

	return result
}

/*
Send sends an email
*/
func (s SmtpMailService) Send(mail Mail) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mail.From.Email, mail.From.Name)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	for _, p := range mail.To {
		m.SetAddressHeader("To", p.Email, p.Name)
	}

	return s.Dialer.DialAndSend(m)
}
