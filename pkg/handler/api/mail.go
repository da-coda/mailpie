package api

import (
	"github.com/gorilla/mux"
	"mailpie/pkg/store"
	"net/http"
)

func registerMailRoutes(router *mux.Router) {
	router.HandleFunc("/", getMails).Methods("GET")
}

func getMails(writer http.ResponseWriter, _ *http.Request) {
	mails := store.GetMails()
	respondWithJSON(writer, 200, mails)
}
