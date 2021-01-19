package imap

import (
	b "github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/backend/memory"
)

type user struct {
	mailboxes map[string]b.Mailbox
	username  string
}

func NewUser(username string) b.User {
	mailboxes := make(map[string]b.Mailbox)
	user := &user{username: username, mailboxes: mailboxes}
	_ = user.CreateMailbox("INBOX")
	_ = user.CreateMailbox("Inbox")
	_ = user.CreateMailbox("Sent")
	return user
}

func (u user) Username() string {
	return u.username
}

func (u user) ListMailboxes(_ bool) ([]b.Mailbox, error) {
	mailboxes := make([]b.Mailbox, 0, len(u.mailboxes))
	for _, mailbox := range u.mailboxes {
		mailboxes = append(mailboxes, mailbox)
	}
	return mailboxes, nil
}

func (u user) GetMailbox(name string) (b.Mailbox, error) {
	mailbox, ok := u.mailboxes[name]
	if !ok {
		return nil, b.ErrNoSuchMailbox
	}
	return mailbox, nil
}

func (u *user) CreateMailbox(name string) error {
	_, ok := u.mailboxes[name]
	if ok {
		return b.ErrMailboxAlreadyExists
	}
	u.mailboxes[name] = &memory.Mailbox{
		Subscribed: true,
		Messages:   nil,
	}
	return nil
}

func (u *user) DeleteMailbox(name string) error {
	delete(u.mailboxes, name)
	return nil
}

func (u *user) RenameMailbox(existingName, newName string) error {
	mailbox, ok := u.mailboxes[existingName]
	if !ok {
		return b.ErrNoSuchMailbox
	}
	_, ok = u.mailboxes[newName]
	if ok {
		return b.ErrMailboxAlreadyExists
	}
	u.mailboxes[newName] = mailbox
	delete(u.mailboxes, existingName)
	return nil
}

func (u user) Logout() error {
	return nil
}
