package api

import (
	"github.com/da-coda/mailpie/pkg/store"
	"github.com/gorilla/mux"
)

var mailStore *store.MailStore

type BaseRouter struct {
	Mux *mux.Router
}

func NewBaseRouter(store *store.MailStore) BaseRouter {
	br := BaseRouter{Mux: mux.NewRouter()}
	mailStore = store
	subrouterV1 := br.Mux.PathPrefix("/v1").Subrouter()
	br.SetupSubrouters(subrouterV1)
	return br
}

func (br *BaseRouter) SetupSubrouters(mux *mux.Router) {
	subrouterMail := mux.PathPrefix("/mail").Subrouter()
	mailSubrouter(subrouterMail)
}
