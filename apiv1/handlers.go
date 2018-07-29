package apiv1

import (
	"net/http"

	"darlinggo.co/api"
	yall "yall.in"
)

func (a APIv1) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	log := yall.FromContext(r.Context())
	accounts, err := a.Accounts.ListAccounts(r.Context())
	if err != nil {
		log.WithError(err).Error("error listing accounts")
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	accts := accountsFromCore(accounts)
	for pos, acct := range accts {
		bal, err := a.Transactions.Balance(r.Context(), acct.ID)
		if err != nil {
			log.WithError(err).WithField("account", acct.ID).Error("error retrieving balance")
			api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
			return
		}
		acct.Balance = bal
		accts[pos] = acct
	}
	api.Encode(w, r, http.StatusOK, Response{Accounts: accts})
}

func (a APIv1) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	api.Encode(w, r, http.StatusOK, Response{})
}

func (a APIv1) handleListGroups(w http.ResponseWriter, r *http.Request) {
	log := yall.FromContext(r.Context())
	groups, err := a.Groups.ListGroups(r.Context())
	if err != nil {
		log.WithError(err).Error("error listing transaction groups")
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Groups: groupsFromCore(groups)})
}

func (a APIv1) handleGetGroup(w http.ResponseWriter, r *http.Request) {
	api.Encode(w, r, http.StatusOK, Response{})
}

func (a APIv1) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	log := yall.FromContext(r.Context())
	filters, apiErrs := txnFilterFromQuery(r.URL.Query())
	if len(apiErrs) > 0 {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: apiErrs})
		return
	}
	txns, err := a.Transactions.ListTransactions(r.Context(), filters)
	if err != nil {
		log.WithError(err).Error("error listing transactions")
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Transactions: txnsFromCore(txns)})
}

func (a APIv1) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	api.Encode(w, r, http.StatusOK, Response{})
}
