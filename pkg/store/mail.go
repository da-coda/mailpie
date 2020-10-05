package store

import (
	"bytes"
	"github.com/DusanKasan/parsemail"
	"mailmon-go/pkg/event"
	"net/mail"
)

type Mail struct {
	From        mail.Address
	To          []mail.Address
	Cc          []mail.Address
	Header      mail.Header
	Body        string
	DecodedBody string
	Raw         []byte
}

func (mail Mail) FromAddress() string {
	return mail.From.Address
}

func (mail Mail) FromName() string {
	return mail.From.Name
}

func (mail Mail) ToAddresses() []mail.Address {
	return mail.To
}

func (mail Mail) CcAddresses() []mail.Address {
	return mail.Cc
}

type Mails []parsemail.Email

func (emails *Mails) add(data []byte) (mailInstance Mail, err error) {

	body, err := parsemail.Parse(bytes.NewReader(data))
	if err != nil {
		return Mail{}, err
	}
	*emails = append(*emails, body)
	events := event.NewOrGet()
	events.Dispatch("mailReceived", "MailHandler", body)
	return
}

var mails Mails

func AddMail(data []byte) (mailInstance Mail, err error) {
	return mails.add(data)
}

func GetMails() Mails {
	return mails
}
