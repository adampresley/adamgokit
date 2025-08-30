package email

import "time"

/*
A Config object tells us how to configure our email server connection
*/
type Config struct {
	ApiKey   string
	Domain   string
	Host     string
	Password string
	Port     int
	Timeout  time.Duration
	UserName string
}
