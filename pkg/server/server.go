package server

import (
	"github.com/emersion/go-smtp"
	"time"
)

func NewDefaultServer(backend smtp.Backend) *smtp.Server{
	s := smtp.NewServer(backend)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	s.EnableBINARYMIME = true
	return s
}