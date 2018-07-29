package apiv1

import (
	"net/http"

	"darlinggo.co/api"
	"github.com/carvers/budget"
	yall "yall.in"
)

// APIv1 holds all the information that we want to
// be available for all the functions in the API,
// things like our storers.
type APIv1 struct {
	budget.Dependencies
}

// Response is used to encode JSON responses; it is
// the global response format for all API responses.
type Response struct {
	Accounts     []Account          `json:"accounts,omitempty"`
	Groups       []Group            `json:"groups,omitempty"`
	Transactions []Transaction      `json:"transactions,omitempty"`
	Errors       []api.RequestError `json:"errors,omitempty"`
}

func (a APIv1) contextLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := a.Log.WithRequest(r)
		r = r.WithContext(yall.InContext(r.Context(), log))
		log.WithField("endpoint", r.Header.Get("Trout-Pattern")).Debug("serving request")
		h.ServeHTTP(w, r)
	})
}
