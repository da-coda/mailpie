package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func registerMailRoutes(router *mux.Router) {
	router.HandleFunc("/", getMails).Methods("GET")
}

func getMails(writer http.ResponseWriter, _ *http.Request) {
	respondWithJSON(writer, 200, nil)
}
