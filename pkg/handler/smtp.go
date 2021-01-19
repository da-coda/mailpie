package handler

import (
	"github.com/sirupsen/logrus"
	"mailpie/pkg/store"
	"net"
)

var SmtpHandler = func(remoteAddr net.Addr, from string, to []string, data []byte) {
	_, err := store.AddMail(data)
	if err != nil {
		logrus.WithError(err).Error("Unable to add mail to store in SMTP handler")
	}
}
