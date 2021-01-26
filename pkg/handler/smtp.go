package handler

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"mailpie/pkg/store"
	"net"
	gomail "net/mail"
	"time"
)

type SmtpHandler struct {
	mailStore store.MailStore
}

func CreateSmtpHandler(mailStore store.MailStore) SmtpHandler {
	return SmtpHandler{mailStore: mailStore}
}

func (handler *SmtpHandler) Handle(_ net.Addr, from string, to []string, data []byte) {
	mail, err := gomail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		logrus.WithError(err).Error("Unable to parse mail in SMTP handler")
	}
	date, err := mail.Header.Date()
	if err != nil {
		logrus.WithError(err).Error("Unable to get date from mail in SMTP handler")
	}
	key := "%s_%s_%s"
	dateString := date.Format(time.RFC3339)
	err = handler.mailStore.Add(fmt.Sprintf(key, from, to, dateString), *mail)

	if err != nil {
		logrus.WithError(err).Error("Unable to add mail to store in SMTP handler")
	}
}
