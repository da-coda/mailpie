package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterApiRoutes(router *mux.Router) {
	mailRouter := router.PathPrefix("/mails/").Subrouter()
	registerMailRoutes(mailRouter)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
