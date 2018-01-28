package budget

import "github.com/apex/log"

//go:generate go-bindata -pkg migrations -o migrations/generated.go sql/

type AccountsStorer interface {
	CreateAccount(Account) error
	GetAccount(string) (Account, error)
	UpdateAccount(string, AccountChange) error
	DeleteAccount(string) error
	ListAccounts() ([]Account, error)
}

type AccountsSensitiveDetailsStorer interface {
	StoreAccountSensitiveDetails(string, AccountSensitiveDetails) error
	GetAccountSensitiveDetails(string) (AccountSensitiveDetails, error)
	DeleteAccountSensitiveDetails(string) error
}

type RecurringStorer interface {
	CreateRecurrings([]Recurring) error
	ListRecurrings() ([]Recurring, error)
	UpdateRecurring(string, RecurringChange) error
}

type TransactionsStorer interface {
	ImportTransactions([]Transaction) error
	ListTransactions(TransactionFilters) ([]Transaction, error)
	UpdateTransactions(TransactionFilters, TransactionChange) error
}

type Dependencies struct {
	Log               *log.Logger
	Accounts          AccountsStorer
	AccountsSensitive AccountsSensitiveDetailsStorer
	Recurring         RecurringStorer
	Transactions      TransactionsStorer
}
