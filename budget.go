package budget

import (
	"context"

	yall "yall.in"
)

//go:generate go-bindata -pkg migrations -o migrations/generated.go sql/

type AccountsStorer interface {
	CreateAccount(context.Context, Account) error
	GetAccount(context.Context, string) (Account, error)
	UpdateAccount(context.Context, string, AccountChange) error
	DeleteAccount(context.Context, string) error
	ListAccounts(context.Context) ([]Account, error)
}

type AccountsSensitiveDetailsStorer interface {
	StoreAccountSensitiveDetails(context.Context, string, AccountSensitiveDetails) error
	GetAccountSensitiveDetails(context.Context, string) (AccountSensitiveDetails, error)
	DeleteAccountSensitiveDetails(context.Context, string) error
}

type RecurringStorer interface {
	CreateRecurrings(context.Context, []Recurring) error
	ListRecurrings(context.Context) ([]Recurring, error)
	UpdateRecurring(context.Context, string, RecurringChange) error
}

type TransactionsStorer interface {
	ImportTransactions(context.Context, []Transaction) error
	ListTransactions(context.Context, TransactionFilters) ([]Transaction, error)
	UpdateTransactions(context.Context, TransactionFilters, TransactionChange) error
	Balance(context.Context, string) (int64, error)
}

type Dependencies struct {
	Log               *yall.Logger
	Accounts          AccountsStorer
	AccountsSensitive AccountsSensitiveDetailsStorer
	Recurring         RecurringStorer
	Transactions      TransactionsStorer
}
