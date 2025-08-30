package email

/*
MailServicer provides an interface describing a service for working with email
*/
type MailServicer interface {
	Send(mail Mail) error
}
