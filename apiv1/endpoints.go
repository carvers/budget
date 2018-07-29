package apiv1

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout"
)

func (a APIv1) Server(baseURL string) http.Handler {
	var router trout.Router
	router.SetPrefix(baseURL)
	router.Endpoint("/transactions/{id}").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleGetTransaction))))
	router.Endpoint("/transactions").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleListTransactions))))

	router.Endpoint("/accounts/{id}").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleGetAccount))))
	router.Endpoint("/accounts").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleListAccounts))))

	router.Endpoint("/groups/{id}").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleGetGroup))))
	router.Endpoint("/groups").Methods("GET").Handler(a.contextLogger(api.NegotiateMiddleware(http.HandlerFunc(a.handleListGroups))))

	return router
}
