package store

import (
	"github.com/da-coda/mailpie/pkg/event"
	"github.com/da-coda/mailpie/pkg/instances"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

var rawMail = []byte(`Received: from localhost (localhost [127.0.0.1])
        by localhost (Mailpie) with SMTP
        for <bob@example.com>; Wed, 27 Jan 2021 17:00:48 +0100 (CET)
MIME-Version: 1.0
Date: Wed, 27 Jan 2021 17:00:48 +0100
From: alex@example.com
To: bob@example.com, cora@example.com
Cc: "Dan" <dan@example.com>
Subject: Hello!
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

Hello <b>Bob</b> and <i>Cora</i>!
`)

type MailStoreUnitTest struct {
	suite.Suite
}

type MockMessageQueue struct {
	mock.Mock
	dispatched map[event.Event]map[string]interface{}
}

func (m *MockMessageQueue) Dispatch(dispatchedEvent event.Event, from string, data interface{}) {
	if m.dispatched == nil {
		m.dispatched = make(map[event.Event]map[string]interface{})
	}
	if m.dispatched[dispatchedEvent] == nil {
		m.dispatched[dispatchedEvent] = make(map[string]interface{})
	}
	m.dispatched[dispatchedEvent][from] = data
	m.Called(dispatchedEvent, from, data)
}

func (suite *MailStoreUnitTest) TestCreate() {
	mockDispatcher := new(MockMessageQueue)
	store := CreateMailStore(mockDispatcher)
	assert.NotNil(suite.T(), store.mails)
}

// Also tests Set()
func (suite *MailStoreUnitTest) TestAdd_NotExist_Dispatch() {
	mockDispatcher := new(MockMessageQueue)
	mail, err := instances.ParseMail(rawMail)
	assert.Nil(suite.T(), err, "Unexpected error")
	store := CreateMailStore(mockDispatcher)

	mockDispatcher.On("Dispatch", NewMailStoredEvent, EventDispatcher, *mail).Return()
	err = store.Add("test", *mail)
	assert.Nil(suite.T(), err, "Unexpected error")
	mockDispatcher.AssertCalled(suite.T(), "Dispatch", NewMailStoredEvent, EventDispatcher, *mail)

	assert.Contains(suite.T(), store.mails, "test", "Mail was not correctly added")
	mailFromStore := store.mails["test"]
	readFromStore, err := ioutil.ReadAll(&mailFromStore)
	assert.Equal(suite.T(), rawMail, readFromStore, "Mail in store isn't equal to original mail")
}

func (suite *MailStoreUnitTest) TestAdd_Exist() {
	mockDispatcher := new(MockMessageQueue)
	mail, err := instances.ParseMail(rawMail)
	assert.Nil(suite.T(), err, "Unexpected error")
	store := CreateMailStore(mockDispatcher)
	store.mails["test"] = *mail
	err = store.Add("test", *mail)
	assert.ErrorIs(suite.T(), err, AlreadyExistsError)
}

func (suite *MailStoreUnitTest) TestGetSingle_Exists() {
	mockDispatcher := new(MockMessageQueue)
	mail, err := instances.ParseMail(rawMail)
	assert.Nil(suite.T(), err, "Unexpected error")
	store := CreateMailStore(mockDispatcher)
	store.mails["test"] = *mail
	single, err := store.GetSingle("test")
	assert.Nil(suite.T(), err, "Unexpected error")
	readFromStore, err := ioutil.ReadAll(&single)
	assert.Equal(suite.T(), rawMail, readFromStore, "Mail in store isn't equal to original mail")
}

func (suite *MailStoreUnitTest) TestGetSingle_NotExists() {
	mockDispatcher := new(MockMessageQueue)
	store := CreateMailStore(mockDispatcher)
	single, err := store.GetSingle("test")
	assert.Equal(suite.T(), single, instances.Mail{}, "Mail should not exist")
	assert.ErrorIs(suite.T(), err, KeyNotExistsError)
}

func (suite *MailStoreUnitTest) TestGetMultiple() {
	mockDispatcher := new(MockMessageQueue)
	store := CreateMailStore(mockDispatcher)
	mail, err := instances.ParseMail(rawMail)
	assert.Nil(suite.T(), err, "Unexpected error")
	store.mails["test"] = *mail
	store.mails["othertest"] = *mail
	store.mails["ignoredtest"] = *mail
	multiple, err, keys := store.GetMultiple([]string{"test", "othertest", "nonexistingtest"})
	assert.Contains(suite.T(), multiple, "test", "Key 'test' should exist")
	assert.Contains(suite.T(), multiple, "othertest", "Key 'othertest' should exist")
	assert.NotContains(suite.T(), multiple, "nonexistingtest", "Key 'nonexistingtest' should not exist")
	assert.ErrorIs(suite.T(), err, KeyNotExistsError)
	assert.Contains(suite.T(), keys, "nonexistingtest", "Key 'nonexistingtest' should be in returned error keys")
}

func TestMailStore(t *testing.T) {
	suite.Run(t, new(MailStoreUnitTest))
}
