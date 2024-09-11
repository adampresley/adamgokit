package email

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

type emailDynamicTemplateRequest struct {
	From             EmailAddress            `json:"from"`
	Personalizations []emailPersonalizations `json:"personalizations"`
	TemplateID       string                  `json:"template_id"`
}

type emailPersonalizations struct {
	To           []EmailAddress         `json:"to"`
	TemplateData map[string]interface{} `json:"dynamic_template_data"`
}

type SendGridBuilder interface {
	GetFrom() EmailAddress
	GetTo() []EmailAddress
	GetTemplateData() map[string]map[string]interface{}
	GetTemplateID() string
	From(email, name string) SendGridBuilder
	TemplateData(to string, data map[string]interface{}) SendGridBuilder
	TemplateDataForAll(data map[string]interface{}) SendGridBuilder
	To(email, name string) SendGridBuilder
	ToMultipleAddresses(emails []string) SendGridBuilder
}

type SendGridEmailBuilder struct {
	templateID   string
	from         EmailAddress
	to           []EmailAddress
	templateData map[string]map[string]interface{}
}

type SendGridServicer interface {
	Send(builder SendGridBuilder) error
}

type SendGridServiceConfig struct {
	ApiKey string
}

type sendGridService struct {
	apiKey string
}

func NewSendGridEmailBuilder(templateID string) *SendGridEmailBuilder {
	return &SendGridEmailBuilder{
		templateID:   templateID,
		from:         EmailAddress{},
		to:           []EmailAddress{},
		templateData: map[string]map[string]interface{}{},
	}
}

func (s *SendGridEmailBuilder) From(email, name string) SendGridBuilder {
	s.from = EmailAddress{
		Email: email,
		Name:  name,
	}

	return s
}

func (s *SendGridEmailBuilder) TemplateData(to string, data map[string]interface{}) SendGridBuilder {
	s.templateData[to] = data
	return s
}

func (s *SendGridEmailBuilder) TemplateDataForAll(data map[string]interface{}) SendGridBuilder {
	for _, to := range s.to {
		s.templateData[to.Email] = data
	}

	return s
}

func (s *SendGridEmailBuilder) To(email, name string) SendGridBuilder {
	s.to = append(s.to, EmailAddress{
		Email: email,
		Name:  name,
	})

	return s
}

func (s *SendGridEmailBuilder) ToMultipleAddresses(emails []string) SendGridBuilder {
	for _, email := range emails {
		s.to = append(s.to, EmailAddress{
			Email: email,
		})
	}

	return s
}

func (s *SendGridEmailBuilder) GetFrom() EmailAddress {
	return s.from
}

func (s *SendGridEmailBuilder) GetTo() []EmailAddress {
	return s.to
}

func (s *SendGridEmailBuilder) GetTemplateData() map[string]map[string]interface{} {
	return s.templateData
}

func (s *SendGridEmailBuilder) GetTemplateID() string {
	return s.templateID
}

func NewSendGridService(config SendGridServiceConfig) *sendGridService {
	return &sendGridService{
		apiKey: config.ApiKey,
	}
}

func (s *sendGridService) Send(builder SendGridBuilder) error {
	var (
		err         error
		request     rest.Request
		requestBody []byte
		response    *rest.Response
	)

	from := builder.GetFrom()
	to := builder.GetTo()
	templateData := builder.GetTemplateData()
	templateID := builder.GetTemplateID()

	if from.Email == "" {
		return fmt.Errorf("from email address is required")
	}

	if len(to) == 0 {
		return fmt.Errorf("to email address is required")
	}

	templateRequest := emailDynamicTemplateRequest{
		From:             from,
		Personalizations: []emailPersonalizations{},
		TemplateID:       templateID,
	}

	if len(templateData) > 0 {
		for to, data := range templateData {
			templateRequest.Personalizations = append(templateRequest.Personalizations, emailPersonalizations{
				To: []EmailAddress{
					{
						Email: to,
					},
				},
				TemplateData: data,
			})
		}
	}

	if requestBody, err = json.Marshal(templateRequest); err != nil {
		return fmt.Errorf("error parsing Send Grid dynamic template request body: %w", err)
	}

	slog.Debug("sending email", slog.Any("request", string(requestBody)))

	request = sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
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
