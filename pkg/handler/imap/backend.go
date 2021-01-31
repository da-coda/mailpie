package imap

import (
	"github.com/emersion/go-imap"
	imapBackend "github.com/emersion/go-imap/backend"
	"github.com/sirupsen/logrus"
	"mailpie/pkg/event"
	"mailpie/pkg/instances"
	"mailpie/pkg/store"
	"time"
)

type backend struct {
	Magpie        imapBackend.User
	UpdateChannel chan imapBackend.Update
}

func NewBackend(mailStore store.MailStore) imapBackend.Backend {
	user := NewUser("Magpie", mailStore)
	updates := make(chan imapBackend.Update)
	backend := &backend{Magpie: user, UpdateChannel: updates}
	events := event.CreateOrGet()
	events.Subscribe(store.NewMailStoredEvent, backend.Handler)
	return backend
}

func (b backend) Updates() <-chan imapBackend.Update {
	return b.UpdateChannel
}

func (b backend) Login(_ *imap.ConnInfo, _, _ string) (imapBackend.User, error) {
	return b.Magpie, nil
}

func (b backend) Handler(_ string, data interface{}) {
	mail := data.(instances.Mail)
	mb, err := b.Magpie.GetMailbox("INBOX")

	if err != nil {
		logrus.WithError(err).Error("Unable to get Mailbox 'INBOX' in IMAP handler")
	}

	err = mb.CreateMessage([]string{imap.RecentFlag}, time.Now(), &mail)
	if err != nil {
		logrus.WithError(err).Error("Unable to create message in IMAP handler")
	}
	update := imapBackend.NewUpdate("Magpie", "INBOX")
	mailboxStatus, err := mb.Status([]imap.StatusItem{imap.StatusMessages, imap.StatusUidNext, imap.StatusRecent, imap.StatusUidValidity, imap.StatusRecent})
	if err != nil {
		logrus.WithError(err).Error("Unable to get Mailbox status")
	}
	b.UpdateChannel <- &imapBackend.MailboxUpdate{
		Update:        update,
		MailboxStatus: mailboxStatus,
	}
}
