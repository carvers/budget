package budget

import (
	"sort"
	"time"
)

type Transaction struct {
	// UUID for the transaction, assigned when it's imported
	ID string

	// ID of the Account this Transaction belongs to
	AccountID string

	// Type of transaction
	TransactionType string
	// Date the transaction was posted
	DatePosted time.Time
	// Date the user initiated the transaction, if known
	DateUser time.Time
	// Date the funds are available
	DateAvailable time.Time
	// Amount of the transaction, in whole cents
	Amount int64
	// The ID for the transaction, assigned by the financial institution
	FiTID string
	// Server-assigned transaction ID, for transactions initiated by the server
	ServerTID string
	// Check (or other reference) number
	CheckNum string
	// Reference number to uniquely identify this transaction. Alternative or addition to CheckNum.
	RefNum string
	// Standard Industrial Code
	SIC int64
	// Payee identifier, if available
	PayeeID string
	// Name of payee or description of transaction
	Name string
	// Extended name of payee or description of transaction
	ExtendedName string
	// Extra information not in Name
	Memo string
	// Source of cash for this transaction, if a 401k transaction.
	Inv401kSource string `sql_column:"inv_401k_source"`
	// ID of the Recurring group to associate this transaction with.
	RecurringID string

	// The three letter code for the currency the amount is specified in
	Currency string
	// The three letter code for the currency the transaction was in before being converted
	OriginalCurrency string

	// The FiTID of a previously downloaded transaction that is being corrected by this transaction
	CorrectFiTID string
	// REPLACE or DELETE; REPLACE means this transaction should replace CorrectFiTID; delete means just remove it.
	CorrectAction string

	// Name of the Payee, if available
	PayeeName string
	// First line of the Payee's address, if available
	PayeeAddr1 string
	// Second line of the Payee's address, if available
	PayeeAddr2 string
	// Third line of the Payee's address, if available
	PayeeAddr3 string
	// City the payee is in, if available
	PayeeCity string
	// State the payee is in, if available
	PayeeState string
	// Postal code the payee is in, if available
	PayeePostalCode string
	// Country the payee is in, if available
	PayeeCountry string
	// Phone number for the payee, if available
	PayeePhone string

	// If this was a transfer to a bank account, and that account's information is available, this is the
	// BankID for that account.
	BankAccountToBankID string
	// If this was a transfer to a bank account, and that account's information is available, this is the
	// BranchID for that account.
	BankAccountToBranchID string
	// If this was a transfer to a bank account, and that account's information is available, this is the
	// AccountID for that account.
	BankAccountToAccountID string
	// If this was a transfer to a bank account, and that account's information is available, this is the
	// AccountType for that account.
	BankAccountToAccountType string
	// If this was a transfer to a bank account, and that account's information is available, this is the
	// AccountKey for that account.
	BankAccountToAccountKey string

	// If this was a transfer to a credit card account, and that account's information is available, this
	// is the AccountID for that account.
	CreditCardAccountToAccountID string
	// If this was a transfer to a credit card account, and that account's information is available, this
	// is the AccountKey for that account.
	CreditCardAccountToAccountKey string
}

func (t Transaction) GetSQLTableName() string {
	return "ofx_transactions"
}

func TransactionsByDate(t []Transaction) {
	sort.Slice(t, func(i, j int) bool {
		return t[i].DatePosted.Before(t[j].DatePosted)
	})
}

type TransactionFilters struct {
	IDs               []string
	AccountIDs        []string
	TransactionTypes  []string
	DatePostedBefore  *time.Time
	DatePostedAfter   *time.Time
	Amount            *int64
	AmountGreaterThan *int64
	AmountLessThan    *int64
	CheckNum          *string
	RefNum            *string
	Name              *string
	RecurringID       *string
}

type TransactionChange struct {
	RecurringID *string
}

func (c TransactionChange) IsEmpty() bool {
	if c.RecurringID != nil {
		return false
	}
	return true
}
