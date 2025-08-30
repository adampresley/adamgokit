package email

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

type sendGridEmailDynamicTemplateRequest struct {
	From             EmailAddress                    `json:"from"`
	Personalizations []sendGridEmailPersonalizations `json:"personalizations"`
	TemplateID       string                          `json:"template_id"`
}

type sendGridEmailPersonalizations struct {
	To           []EmailAddress `json:"to"`
	TemplateData map[string]any `json:"dynamic_template_data"`
}

type SendGridService struct {
	Config *Config
}

func NewSendGridService(config *Config) *SendGridService {
	return &SendGridService{
		Config: config,
	}
}

func (s *SendGridService) Send(mail Mail) error {
	var (
		err         error
		request     rest.Request
		requestBody []byte
		response    *rest.Response
	)

	if mail.From.Email == "" {
		return fmt.Errorf("from email address is required")
	}

	if len(mail.To) == 0 {
		return fmt.Errorf("to email address is required")
	}

	templateRequest := sendGridEmailDynamicTemplateRequest{
		From:             mail.From,
		Personalizations: []sendGridEmailPersonalizations{},
		TemplateID:       mail.Template,
	}

	if len(mail.TemplateData) > 0 {
		templateRequest.Personalizations = append(templateRequest.Personalizations, sendGridEmailPersonalizations{
			To:           mail.To,
			TemplateData: mail.TemplateData,
		})
	}

	if requestBody, err = json.Marshal(templateRequest); err != nil {
		return fmt.Errorf("error parsing Send Grid dynamic template request body: %w", err)
	}

	slog.Debug("sending email", slog.Any("request", string(requestBody)))

	request = sendgrid.GetRequest(s.Config.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = http.MethodPost
	request.Body = requestBody

	if response, err = sendgrid.API(request); err != nil {
		return fmt.Errorf("error sending mail request to SendGrid: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 399 {
		return fmt.Errorf("error sending mail: %s", response.Body)
	}

	return nil
}
