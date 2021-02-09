package store

import (
	"mailpie/pkg/event"
	"mailpie/pkg/instances"
)

const NewMailStoredEvent event.Event = "newMailStored"
const EventDispatcher = "MailStore"

//MailStore holds a bunch of instances.Mail within a map and notifies via the message queue on mail updates
type MailStore struct {
	mails        map[string]instances.Mail
	messageQueue event.Dispatcher
}

//CreateMailStore always creates and returns a new MailStore. Needs a event.Dispatcher to notify others on mail updates
func CreateMailStore(messageQueue event.Dispatcher) *MailStore {
	var store *MailStore
	store = &MailStore{messageQueue: messageQueue}
	store.mails = make(map[string]instances.Mail)
	return store
}

//Add puts a instances.Mail into the internal map with the given key if the key not exists. Key can be any string but should be recreatable for receiving purposes
//Returns an AlreadyExistsError if key exists in Map
func (store *MailStore) Add(key string, mailData instances.Mail) error {
	_, exists := store.mails[key]
	if exists {
		return AlreadyExistsError
	}

	store.Set(key, mailData)
	return nil
}

//Set puts a instances.Mail into the internal map with the given key, regardless of key existence
func (store *MailStore) Set(key string, data instances.Mail) {
	store.mails[key] = data
	store.messageQueue.Dispatch(NewMailStoredEvent, EventDispatcher, data)
}

// GetSingle for retrieving a single mail by key.
// Returns KeyNotExistsError if given key does not exist in internal map.
func (store *MailStore) GetSingle(key string) (instances.Mail, error) {
	mail, exists := store.mails[key]
	if !exists {
		return instances.Mail{}, KeyNotExistsError
	}
	return mail, nil
}

// GetMultiple retrieves multiple mails for a given key slice. If any of the keys not exist, a KeyNotExistsError will be returned
// but the function will still gather the rest. Any not found key will be within the notFoundKeys return parameter
func (store *MailStore) GetMultiple(keys []string) (mails map[string]instances.Mail, err error, notFoundKeys []string) {
	mails = make(map[string]instances.Mail)
	for _, key := range keys {
		mail, returnedErr := store.GetSingle(key)
		if returnedErr == KeyNotExistsError {
			err = KeyNotExistsError
			notFoundKeys = append(notFoundKeys, key)
			continue
		}
		mails[key] = mail
	}
	return
}
