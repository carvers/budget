package budget

import "github.com/pkg/errors"

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Account struct {
	ID          string
	Name        string
	AccountType string
	RequestType string
	BankORG     string
	BankFID     string
	BankURL     string
}

func (a Account) GetSQLTableName() string {
	return "accounts"
}

type AccountChange struct {
	Name        *string
	AccountType *string
	RequestType *string
	BankORG     *string
	BankFID     *string
	BankURL     *string
}

func (a AccountChange) IsEmpty() bool {
	if a.Name != nil {
		return false
	}
	if a.AccountType != nil {
		return false
	}
	if a.RequestType != nil {
		return false
	}
	if a.BankORG != nil {
		return false
	}
	if a.BankFID != nil {
		return false
	}
	if a.BankURL != nil {
		return false
	}
	return true
}

type AccountSensitiveDetails struct {
	AccountID string
	BankID    string
	UserID    string
	UserPass  string
}
