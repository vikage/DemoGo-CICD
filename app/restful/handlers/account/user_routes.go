package account

import (
	"github.com/gorilla/mux"
)

// RegisterHandlerForAccountAPI register handle use for account APIs
func RegisterHandlerForAccountAPI(r *mux.Router, preAuthenticateRouter *mux.Router) {
	r.HandleFunc("/login", LoginRequestHandler).Methods("POST")
}
