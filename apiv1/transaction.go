package apiv1

import (
	"net/url"
	"strconv"
	"time"

	"darlinggo.co/api"
	"github.com/carvers/budget"
)

type Transaction struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	TransactionType string    `json:"transaction_type"`
	DatePosted      time.Time `json:"date_posted"`
	DateUser        time.Time `json:"date_user,omitempty"`
	DateAvailable   time.Time `json:"date_available,omitempty"`
	Amount          int64     `json:"amount"`
	Name            string    `json:"name"`
	ExtendedName    string    `json:"extended_name"`
	Memo            string    `json:"memo"`
	EditedName      string    `json:"edited_name"`
	GroupID         string    `json:"group_id"`
}

func txnFromCore(txn budget.Transaction) Transaction {
	return Transaction{
		ID:              txn.ID,
		AccountID:       txn.AccountID,
		TransactionType: txn.TransactionType,
		DatePosted:      txn.DatePosted,
		DateUser:        txn.DateUser,
		DateAvailable:   txn.DateAvailable,
		Amount:          txn.Amount,
		Name:            txn.Name,
		ExtendedName:    txn.ExtendedName,
		Memo:            txn.Memo,
		EditedName:      txn.EditedName,
		GroupID:         txn.GroupID,
	}
}

func txnsFromCore(txns []budget.Transaction) []Transaction {
	res := make([]Transaction, 0, len(txns))
	for _, txn := range txns {
		res = append(res, txnFromCore(txn))
	}
	return res
}

func txnFilterFromQuery(q url.Values) (budget.TransactionFilters, []api.RequestError) {
	tf := budget.TransactionFilters{
		IDs:              q["id"],
		AccountIDs:       q["account_id"],
		TransactionTypes: q["transaction_type"],
	}
	var errs []api.RequestError
	if q.Get("date_posted_before") != "" {
		t, err := time.Parse(time.RFC3339, q.Get("date_posted_before"))
		if err != nil {
			errs = append(errs, api.RequestError{
				Slug:  api.RequestErrInvalidFormat,
				Param: "date_posted_before",
			})
		} else {
			tf.DatePostedBefore = &t
		}
	}
	if q.Get("date_posted_after") != "" {
		t, err := time.Parse(time.RFC3339, q.Get("date_posted_after"))
		if err != nil {
			errs = append(errs, api.RequestError{
				Slug:  api.RequestErrInvalidFormat,
				Param: "date_posted_after",
			})
		} else {
			tf.DatePostedAfter = &t
		}
	}
	if q.Get("amount") != "" {
		amt, err := strconv.ParseInt(q.Get("amount"), 10, 64)
		if err != nil {
			errs = append(errs, api.RequestError{
				Slug:  api.RequestErrInvalidFormat,
				Param: "amount",
			})
		} else {
			tf.Amount = &amt
		}
	}
	if q.Get("amount_greater_than") != "" {
		amt, err := strconv.ParseInt(q.Get("amount_greater_than"), 10, 64)
		if err != nil {
			errs = append(errs, api.RequestError{
				Slug:  api.RequestErrInvalidFormat,
				Param: "amount_greater_than",
			})
		}
		tf.AmountGreaterThan = &amt
	}
	if q.Get("amount_less_than") != "" {
		amt, err := strconv.ParseInt(q.Get("amount_less_than"), 10, 64)
		if err != nil {
			errs = append(errs, api.RequestError{
				Slug:  api.RequestErrInvalidFormat,
				Param: "amount_less_than",
			})
		}
		tf.AmountLessThan = &amt
	}
	if q.Get("check_num") != "" {
		num := q.Get("check_num")
		tf.CheckNum = &num
	}
	if q.Get("ref_num") != "" {
		num := q.Get("ref_num")
		tf.RefNum = &num
	}
	if q.Get("name") != "" {
		name := q.Get("name")
		tf.Name = &name
	}
	if q.Get("group_id") != "" {
		id := q.Get("group_id")
		tf.GroupID = &id
	}
	return tf, errs
}
