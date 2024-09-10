package auth

import "github.com/adampresley/goth"

type DirectUserLoginInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	OrgID    string `json:"orgID"`
}

type DirectConfig struct {
	UserValidator func(loginInput DirectUserLoginInput) (goth.User, error)
}
