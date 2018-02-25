package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/carvers/budget"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
)

func transactionFromRow(row []string, accountID string) (budget.Transaction, error) {
	txnID := row[0]
	date := row[1]
	name := row[2]
	amountStr := row[3]
	isInterest := row[4]
	id, err := uuid.GenerateUUID()
	if err != nil {
		return budget.Transaction{}, errors.Wrap(err, "error generating UUID")
	}
	amountStr = strings.Replace(amountStr, ".", "", -1)
	amountStr = strings.Replace(amountStr, ",", "", -1)
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return budget.Transaction{}, errors.Wrap(err, fmt.Sprintf("error parsing amount %q as int", amountStr))
	}
	amount = amount * -1
	transactionType := "DEBIT"
	if strings.TrimSpace(strings.ToLower(isInterest)) == "y" {
		transactionType = "FEE"
	} else if amount > 0 {
		transactionType = "CREDIT"
	}
	dateTime, err := time.ParseInLocation("1/2/2006", date, time.Local)
	if err != nil {
		return budget.Transaction{}, errors.Wrap(err, fmt.Sprintf("error parsing date %q as time.Time", date))
	}
	return budget.Transaction{
		ID:              id,
		AccountID:       accountID,
		Amount:          amount,
		FiTID:           txnID,
		Name:            name,
		DatePosted:      dateTime,
		TransactionType: transactionType,
	}, nil
}

func FromReader(r io.Reader, accountID string) ([]budget.Transaction, error) {
	c := csv.NewReader(r)
	rows, err := c.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "error reading CSV")
	}
	if rc, ok := r.(io.ReadCloser); ok {
		rc.Close()
	}
	txns := make([]budget.Transaction, 0, len(rows))
	for _, row := range rows {
		txn, err := transactionFromRow(row, accountID)
		if err != nil {
			return nil, errors.Wrap(err, "error converting csv row")
		}
		txns = append(txns, txn)
	}
	return txns, nil
}
