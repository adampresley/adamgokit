package auth

import (
	"github.com/adampresley/goth/providers/direct"
)

type DirectConfig struct {
	LoginURI    string
	UserFetcher direct.UserFetcher
	CredChecker direct.CredChecker
}
