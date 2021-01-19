package handler

import (
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"mailpie/pkg/event"
	"net/http"
)

type SSEHandler struct {
	Server *sse.Server
}

var sseHandler *SSEHandler

func NewOrGetSSEHandler(Server *sse.Server) *SSEHandler {
	if sseHandler != nil {
		return sseHandler
	}
	sseHandler = &SSEHandler{Server}
	events := event.NewOrGet()
	events.Subscribe("mailReceived", sseHandler.Publish)
	return sseHandler
}

func (sseHandler *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sseHandler.Server.HTTPHandler(w, r)
}

func (sseHandler *SSEHandler) Publish(from string, data interface{}) {
	payload, err := json.Marshal(data.(parsemail.Email))
	if err != nil {
		logrus.WithError(err).Error("Unable to Marshal email in SSEHandler")
		fmt.Println(err.Error())
	}
	e := sse.Event{ID: []byte("1"), Data: payload, Event: []byte("message")}
	sseHandler.Server.Publish("messages", &e)
}
