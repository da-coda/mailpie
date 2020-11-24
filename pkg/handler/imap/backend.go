package imap

import (
	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-imap"
	b "github.com/emersion/go-imap/backend"
	"io"
	"mailpie/pkg/event"
	"time"
)

type backend struct {
	Magpie        b.User
	UpdateChannel chan b.Update
}

func NewBackend() b.Backend {
	user := NewUser("Magpie")
	updates := make(chan b.Update)
	backend := &backend{Magpie: user, UpdateChannel: updates}
	events := event.NewOrGet()
	events.Subscribe("mailReceived", backend.Handler)
	return backend
}

func (b backend) Updates() <-chan b.Update {
	return b.UpdateChannel
}

func (b backend) Login(_ *imap.ConnInfo, _, _ string) (b.User, error) {
	return b.Magpie, nil
}

func (b backend) Handler(dispatcher string, data interface{}) {
	mail := data.(parsemail.Email)
	wrappedMail := WrappedParsemail{mail.Content, mail}
	mb, _ := b.Magpie.GetMailbox("Inbox")
	mb.CreateMessage([]string{imap.RecentFlag}, time.Now(), wrappedMail)
}

type WrappedParsemail struct {
	io.Reader
	parsemail.Email
}

func (w WrappedParsemail) Len() int {
	var bytes []byte
	_, _ = w.Content.Read(bytes)
	return len(bytes)
}
