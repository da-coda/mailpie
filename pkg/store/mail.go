package store

import (
	"mailpie/pkg/event"
	gomail "net/mail"
)

const NewMailStoredEvent event.Event = "newMailStored"
const EventDispatcher = "MailStore"

type MailStore struct {
	mails        map[string]gomail.Message
	messageQueue event.MessageQueue
}

var mailStore *MailStore

func CreateMailStore(messageQueue event.MessageQueue) *MailStore {
	var store *MailStore
	store = &MailStore{messageQueue: messageQueue}
	store.mails = make(map[string]gomail.Message)
	if mailStore == nil {
		mailStore = store
	}
	return mailStore
}

func (store *MailStore) Add(key string, mailData gomail.Message) error {
	_, exists := store.mails[key]
	if exists {
		return AlreadyExistsError
	}

	return store.Set(key, mailData)
}

func (store *MailStore) Set(key string, data gomail.Message) error {
	store.mails[key] = data
	store.messageQueue.Dispatch(NewMailStoredEvent, EventDispatcher, data)
	return nil
}

func (store *MailStore) Get() (map[string]gomail.Message, error) {
	return store.mails, nil
}

func (store *MailStore) GetSingle(key string) (gomail.Message, error) {
	mail, exists := store.mails[key]
	if !exists {
		return gomail.Message{}, KeyNotExistsError
	}
	return mail, nil
}

func (store *MailStore) GetMultiple(keys []string) (mails []gomail.Message, err error, notFoundKeys []string) {
	for _, key := range keys {
		mail, returnedErr := store.GetSingle(key)
		if returnedErr == KeyNotExistsError {
			err = KeyNotExistsError
			notFoundKeys = append(notFoundKeys, key)
		}
		mails = append(mails, mail)
	}
	return
}
