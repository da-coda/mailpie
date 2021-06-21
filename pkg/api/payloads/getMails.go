package payloads

import (
	"fmt"
	"github.com/da-coda/mailpie/pkg/instances"
	"github.com/pkg/errors"
)

type GetMailsItem struct {
	Id      string   `json:"id"`
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
}

type GetMailsPayload []GetMailsItem

func (mailItem *GetMailsItem) FromGoMail(mail instances.Mail, id string) error {
	mailItem.From = mail.Header.Get("From")
	mailItem.Subject = mail.Header.Get("Subject")
	mailItem.Id = id
	toAddresses, err := mail.Header.AddressList("To")
	if err != nil {
		return errors.Wrap(err, "unable to get to addresses")
	}
	for _, address := range toAddresses {
		mailItem.To = append(mailItem.To, fmt.Sprintf("%s <%s>", address.Name, address.Address))
	}
	return nil
}
