package apiv1

import "github.com/carvers/budget"

type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

func accountFromCore(a budget.Account) Account {
	return Account{
		ID:   a.ID,
		Name: a.Name,
	}
}

func accountsFromCore(as []budget.Account) []Account {
	res := make([]Account, 0, len(as))
	for _, a := range as {
		res = append(res, accountFromCore(a))
	}
	return res
}
