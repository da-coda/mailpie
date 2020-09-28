package api

import (
	"github.com/gorilla/mux"
	"mailmon-go/pkg/store"
	"net/http"
)

func registerMailRoutes(router *mux.Router) {
	router.HandleFunc("/", getMails).Methods("GET")
}

func getMails(writer http.ResponseWriter, request *http.Request) {
	mails := store.GetMails()
	respondWithJSON(writer, 200, mails)
}
