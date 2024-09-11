package email

/*
Mail represents an email. Who's sending, recipients, subject, and message
*/
type Mail struct {
	Body    string
	From    EmailAddress
	Subject string
	To      []EmailAddress
}
