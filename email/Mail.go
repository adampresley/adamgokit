package email

/*
Mail represents an email. Who's sending, recipients, subject, and message.
It is a basic structure for sending email for most systems.
*/
type Mail struct {
	Body         string
	BodyIsHtml   bool
	From         EmailAddress
	Subject      string
	Template     string
	TemplateData map[string]any
	To           []EmailAddress
}
