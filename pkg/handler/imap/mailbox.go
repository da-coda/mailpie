package imap

import (
	"github.com/emersion/go-imap"
	"mailpie/pkg/store"
	"strings"
	"time"
)

type Mailbox struct {
	reference string
	Mailstore *store.MailStore
	name      string
	User      *user
}

func (m *Mailbox) SetSubscribed(subscribed bool) error {
	return nil
}

func (m *Mailbox) Check() error {
	panic("implement me")
}

func (m *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	panic("implement me")
}

func (m *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	panic("implement me")
}

func (m *Mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	panic("implement me")
}

func (m *Mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	panic("implement me")
}

func (m *Mailbox) Expunge() error {
	panic("implement me")
}

func (m *Mailbox) Name() string {
	return strings.TrimPrefix(m.name, m.reference)
}

func (m *Mailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Name: m.name,
	}
	return info, nil
}

func (m *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(m.name, items)
	status.Flags = m.flags()
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = m.unseenSeqNum()

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = uint32(len(m.Mailstore.GetAll()))
		case imap.StatusUidNext:
			status.UidNext = m.uidNext()
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = 0 // TODO
		case imap.StatusUnseen:
			status.Unseen = 0 // TODO
		}
	}

	return status, nil
}

func (m *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	for i, msg := range m.Mailstore.GetAll() {
		seqNum, _ := m.Mailstore.GetSeqNumberForKey(i)

		var id uint32
		if uid {
			id = msg.Get32BitTimestamp()
		} else {
			id = seqNum
		}
		if !seqSet.Contains(id) {
			continue
		}

		ch <- msg
	}

	return nil
}

func (m *Mailbox) uidNext() uint32 {
	var uid uint32
	for _, msg := range m.Mailstore.GetAll() {
		uidTemp := msg.Get32BitTimestamp()
		if msg.Get32BitTimestamp() > uid {
			uid = uidTemp
		}
	}
	uid++
	return uid
}

func (m *Mailbox) flags() []string {
	flagsMap := make(map[string]bool)
	for _, msg := range m.Mailstore.GetAll() {
		for _, f := range msg.Flags {
			if !flagsMap[f] {
				flagsMap[f] = true
			}
		}
	}

	var flags []string
	for f := range flagsMap {
		flags = append(flags, f)
	}
	return flags
}
