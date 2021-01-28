package event

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"runtime"
	"testing"
)

type EventsTestSuite struct {
	suite.Suite
}

func (suite *EventsTestSuite) AfterTest(suiteName, testName string) {
	mq = nil
}

func (suite *EventsTestSuite) TestNewOrGet_New() {
	assert.Nil(suite.T(), mq, "MessageQueue should be nil at the beginning")
	messagequeue := NewOrGet()
	assert.NotNil(suite.T(), messagequeue, "NewOrGet should return an instance, Nil returned")
	assert.NotNil(suite.T(), mq, "mq should not be Nil")
}

func (suite *EventsTestSuite) TestNewOrGet_Get() {
	assert.Nil(suite.T(), mq, "MessageQueue should be nil at the beginning")
	messagequeueNew := NewOrGet()
	messagequeueGet := NewOrGet()
	assert.Same(suite.T(), messagequeueNew, messagequeueGet, "Second NewOrGet call should return the same Instance as the first call")
}

func (suite *EventsTestSuite) TestSubscribe() {
	handler := Handler(func(dispatcher string, data interface{}) { return })
	messagequeue := NewOrGet()
	_, exists := mq.topics["test"]
	assert.False(suite.T(), exists, "Topic 'test' should not exist already")
	messagequeue.Subscribe("test", handler)
	handlersForTopic, exists := mq.topics["test"]
	handlerInTopic := handlersForTopic[0]
	assert.True(suite.T(), exists, "Topic 'test' should exist")
	funcName1 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(handlerInTopic).Pointer()).Name()
	assert.Equal(suite.T(), funcName1, funcName2, "Subscribed handler is not the same as the one found in topic")
}

func (suite *EventsTestSuite) TestDispatch_WithSubscriber() {
	messagequeue := NewOrGet()
	handler := Handler(func(dispatcher string, data interface{}) {
		assert.Equal(suite.T(), "TestDispatch", dispatcher)
		assert.Equal(suite.T(), "This is a Test", data)
	})
	messagequeue.Subscribe("test", handler)
	messagequeue.Dispatch("test", "TestDispatch", "This is a Test")
}

func (suite *EventsTestSuite) TestDispatch_WithoutSubscriber() {
	messagequeue := NewOrGet()
	assert.NotPanics(suite.T(), func() {
		messagequeue.Dispatch("test", "TestDispatch", "This is a Test")
	})
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}