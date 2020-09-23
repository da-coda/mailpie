package backend

import (
	"github.com/emersion/go-smtp"
	"mailmon-go/pkg/session"
)

func New() smtp.Backend {
	return &backend{}
}

// The Backend implements SMTP server methods.
type backend struct{}

// Login handles a login command with username and password.
func (bkd *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &session.Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &session.Session{}, nil
}
