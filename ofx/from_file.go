package ofx

import (
	"io"

	"github.com/aclindsa/ofxgo"
	"github.com/carvers/budget"
	"github.com/pkg/errors"
)

func FromReader(r io.Reader) (budget.AccountSensitiveDetails, []budget.Transaction, error) {
	var asd budget.AccountSensitiveDetails
	resp, err := ofxgo.ParseResponse(r)
	if err != nil {
		return asd, nil, errors.Wrap(err, "error parsing file")
	}
	if rc, ok := r.(io.ReadCloser); ok {
		rc.Close()
	}
	var transactions []budget.Transaction
	for _, msg := range resp.Bank {
		if stmt, ok := msg.(*ofxgo.StatementResponse); ok {
			asd = budget.AccountSensitiveDetails{
				AccountID: string(stmt.BankAcctFrom.AcctID),
				BankID:    string(stmt.BankAcctFrom.BankID),
			}
			txns, err := transactionsFromStatement(stmt)
			if err != nil {
				return asd, nil, errors.Wrap(err, "error parsing transaction")
			}
			transactions = append(transactions, txns...)
		} else {
			return asd, nil, errors.Errorf("unknown message type %T", msg)
		}
	}
	for _, msg := range resp.CreditCard {
		if stmt, ok := msg.(*ofxgo.CCStatementResponse); ok {
			asd = budget.AccountSensitiveDetails{
				AccountID: string(stmt.CCAcctFrom.AcctID),
			}
			txns, err := transactionsFromCCStatement(stmt)
			if err != nil {
				return asd, nil, errors.Wrap(err, "error parsing transaction")
			}
			transactions = append(transactions, txns...)
		} else {
			return asd, nil, errors.Errorf("unknown message type %T", msg)
		}
	}
	return asd, transactions, nil
}
