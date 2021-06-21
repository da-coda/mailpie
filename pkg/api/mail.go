package api

import (
	"encoding/json"
	"github.com/da-coda/mailpie/pkg/api/payloads"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func mailSubrouter(mux *mux.Router) {
	mux.HandleFunc("", getMails).Methods(http.MethodGet)
	mux.HandleFunc("/{key}", getMails).Methods(http.MethodGet)
}

func getMails(w http.ResponseWriter, r *http.Request) {
	payload := payloads.GetMailsPayload{}
	mails := mailStore.GetAll()
	for id, mail := range mails {
		mailItem := payloads.GetMailsItem{}
		err := mailItem.FromGoMail(mail, id)
		if err != nil {
			logrus.WithError(err).WithField("route", "GET v1/mail").Error("Unable to create mail payload")
			w.WriteHeader(500)
			return
		}
		payload = append(payload, mailItem)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).WithField("route", "GET v1/mail").Error("Unable to marshal payload")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(payloadJSON)
	if err != nil {
		logrus.WithError(err).WithField("route", "GET v1/mail").Error("Unable to write response")
		return
	}
}
