package store

import (
	"bytes"
	"github.com/DusanKasan/parsemail"
	"mailpie/pkg/event"
)

type Mails []parsemail.Email

func (emails *Mails) add(data []byte) (mailInstance parsemail.Email, err error) {

	body, err := parsemail.Parse(bytes.NewReader(data))
	if err != nil {
		return parsemail.Email{}, err
	}
	*emails = append(*emails, body)
	events := event.NewOrGet()
	events.Dispatch("mailReceived", "MailHandler", body)
	events.Dispatch("rawMailReceived", "MailHandler", data)
	return body, nil
}

var mails Mails

func AddMail(data []byte) (mailInstance parsemail.Email, err error) {
	return mails.add(data)
}

func GetMails() Mails {
	return mails
}
