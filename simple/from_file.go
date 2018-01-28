package simple

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/carvers/budget"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
)

type category struct {
	UUID     string
	Name     string
	Folder   string
	FolderID int64
}

type times struct {
	WhenRecorded      int64  `json:"when_recorded"`
	WhenRecordedLocal string `json:"when_recorded_local"`
	WhenReceived      int64  `json:"when_received"`
	LastModified      int64  `json:"last_modified"`
	LastTXVia         int64  `json:"last_txvia"`
}

type amounts struct {
	Amount   int64 `json:"amount"`
	Cleared  int64 `json:"cleared"`
	Fees     int64 `json:"fees"`
	Cashback int64 `json:"cashback"`
	Base     int64 `json:"base"`
}

type geo struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type transaction struct {
	UUID            string     `json:"uuid"`
	UserID          string     `json:"user_id"`
	RecordType      string     `json:"record_type"`
	TransactionType string     `json:"transaction_type"`
	BookkeepingType string     `json:"bookkeeping_type"`
	IsHold          bool       `json:"is_hold"`
	IsActive        bool       `json:"is_active"`
	RunningBalance  int64      `json:"running_balance"`
	RawDescription  string     `json:"raw_description"`
	Description     string     `json:"description"`
	Memo            string     `json:"memo"`
	Categories      []category `json:"categories"`
	Times           times      `json:"times"`
	Amounts         amounts    `json:"amounts"`
	CorrelationID   string     `json:"correlation_id"`
	AccountID       string     `json:"account_id"`
	InitiatedBy     string     `json:"initiated_by"`
	Partner         string     `json:"partner"`
	Geo             geo        `json:"geo"`
}

type download struct {
	Transactions []transaction
	Offset       int64
}

func transactionFromTransaction(t transaction) (budget.Transaction, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return budget.Transaction{}, errors.Wrap(err, "error generating UUID")
	}
	types := map[string]string{
		"ach":                   "XFER",
		"ach_reversal":          "XFER",
		"ach_reversal_reversal": "XFER",
		"adjustment":            "XFER",
		"argo_debit_reversal":   "CREDIT",
		"atm_withdrawal":        "ATM",
		"balance_sweep":         "XFER",
		"bill_payment":          "CHECK",
		"bill_payment_reversal": "CHECK",
		"c2c":                          "XFER",
		"check_deposit":                "CHECK",
		"courtesy_credit":              "SRVCHG",
		"fee":                          "FEE",
		"interest_credit":              "INT",
		"make_up_credit":               "SRVCHG",
		"migration_interbank_transfer": "XFER",
		"otc_withdrawal":               "CASH",
		"pin_purchase":                 "DEBIT",
		"provisional_credit":           "SRVCHG",
		"purchase":                     "DEBIT",
		"service_charge_refund":        "SRVCHG",
		"shared_transfer":              "XFER",
		"signature_credit":             "CREDIT",
		"signature_purchase":           "DEBIT",
		"signature_return":             "CREDIT",
	}
	transactionType, ok := types[strings.ToLower(t.TransactionType)]
	if !ok {
		return budget.Transaction{}, errors.Errorf("unknown transfer type %q, please open an issue at https://github.com/carvers/budget/issues/new", t.TransactionType)
	}
	var positive int64
	switch strings.ToLower(t.BookkeepingType) {
	case "credit":
		positive = 1
	case "debit":
		positive = -1
	default:
		return budget.Transaction{}, errors.Errorf("unknown transfer type %q, please open an issue at https://github.com/carvers/budget/issues/new", t.BookkeepingType)
	}
	return budget.Transaction{
		ID:              id,
		AccountID:       t.AccountID,
		Amount:          t.Amounts.Amount / 100 * positive,
		FiTID:           t.UUID,
		Name:            t.Description,
		ExtendedName:    t.RawDescription,
		Memo:            t.Memo,
		DateUser:        time.Unix(t.Times.WhenRecorded/1000, 0),
		DatePosted:      time.Unix(t.Times.WhenReceived/1000, 0),
		TransactionType: transactionType,
		PayeeAddr1:      t.Geo.Street,
		PayeeCity:       t.Geo.City,
		PayeeState:      t.Geo.State,
		PayeePostalCode: t.Geo.Zip,
	}, nil
}

func FromReader(r io.Reader) ([]budget.Transaction, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "error reading stream")
	}
	if rc, ok := r.(io.ReadCloser); ok {
		rc.Close()
	}
	var resp download
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling stream")
	}
	txns := make([]budget.Transaction, 0, len(resp.Transactions))
	for _, trans := range resp.Transactions {
		txn, err := transactionFromTransaction(trans)
		if err != nil {
			return nil, errors.Wrap(err, "error converting Simple transaction")
		}
		txns = append(txns, txn)
	}
	return txns, nil
}
