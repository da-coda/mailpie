package handler

import (
	"fmt"
	"github.com/da-coda/mailpie/pkg/instances"
	"github.com/da-coda/mailpie/pkg/store"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type SmtpHandler struct {
	mailStore store.MailStore
}

func CreateSmtpHandler(mailStore store.MailStore) SmtpHandler {
	return SmtpHandler{mailStore: mailStore}
}

//Handle incoming emails. Parses the incoming mail into instances.Mail and then writes the mail into the mailStore with the key format
// from_to_2006-01-02T15:04:05Z07:00
func (handler *SmtpHandler) Handle(_ net.Addr, from string, to []string, data []byte) {
	mail, err := instances.ParseMail(data)
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
