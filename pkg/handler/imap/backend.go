package imap

import (
	"bytes"
	"github.com/emersion/go-imap"
	imapBackend "github.com/emersion/go-imap/backend"
	"github.com/sirupsen/logrus"
	"io"
	"mailpie/pkg/event"
	"time"
)

type backend struct {
	Magpie        imapBackend.User
	UpdateChannel chan imapBackend.Update
}

func NewBackend() imapBackend.Backend {
	user := NewUser("Magpie")
	updates := make(chan imapBackend.Update)
	backend := &backend{Magpie: user, UpdateChannel: updates}
	events := event.NewOrGet()
	events.Subscribe("rawMailReceived", backend.Handler)
	return backend
}

func (b backend) Updates() <-chan imapBackend.Update {
	return b.UpdateChannel
}

func (b backend) Login(info *imap.ConnInfo, user, password string) (imapBackend.User, error) {
	return b.Magpie, nil
}

func (b backend) Handler(_ string, data interface{}) {
	mail := data.([]byte)

	wrappedMail := WrappedParsemail{bytes.NewReader(mail), len(mail)}
	mb, _ := b.Magpie.GetMailbox("INBOX")
	err := mb.CreateMessage([]string{imap.RecentFlag}, time.Now(), wrappedMail)
	if err != nil {
		logrus.WithError(err).Error("Unable to create message in IMAP handler")
	}
	b.UpdateChannel <- imapBackend.NewUpdate("Magpie", "Inbox")
}

type WrappedParsemail struct {
	io.Reader
	size int
}

func (w WrappedParsemail) Len() int {
	return w.size
}
