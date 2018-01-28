package ofx

import (
	"fmt"
	"math/big"

	"github.com/aclindsa/ofxgo"
	"github.com/carvers/budget"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
)

const (
	ReqTypeBank       = "BANK"
	ReqTypeCreditCard = "CC"
)

func FetchTransactions(d budget.Dependencies, account budget.Account, asd budget.AccountSensitiveDetails, clientUID string) ([]budget.Transaction, error) {
	uid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating UUID for transaction")
	}

	query := ofxgo.Request{
		URL: account.BankURL,
		Signon: ofxgo.SignonRequest{
			Org:       ofxgo.String(account.BankORG),
			Fid:       ofxgo.String(account.BankFID),
			UserID:    ofxgo.String(asd.UserID),
			UserPass:  ofxgo.String(asd.UserPass),
			ClientUID: ofxgo.UID(clientUID),
		},
	}

	switch account.RequestType {
	case ReqTypeBank:
		// if bank request..
		accttype, err := ofxgo.NewAcctType(account.AccountType)
		if err != nil {
			return nil, errors.Wrapf(err, "%q is not a valid AcctType", account.AccountType)
		}
		query.Bank = []ofxgo.Message{
			&ofxgo.StatementRequest{
				TrnUID: ofxgo.UID(uid),
				BankAcctFrom: ofxgo.BankAcct{
					BankID:   ofxgo.String(asd.BankID),
					AcctID:   ofxgo.String(asd.AccountID),
					AcctType: accttype,
				},
				Include:        true,
				IncludePending: true,
			},
		}
	case ReqTypeCreditCard:
		// if credit card...
		query.CreditCard = []ofxgo.Message{
			&ofxgo.CCStatementRequest{
				TrnUID: ofxgo.UID(uid),
				CCAcctFrom: ofxgo.CCAcct{
					AcctID: ofxgo.String(asd.AccountID),
				},
				Include:        true,
				IncludePending: true,
			},
		}
	default:
		return nil, errors.Errorf("unknown ReqType %q", account.RequestType)
	}

	client := ofxgo.Client{
		AppID:       "CARVE",
		AppVer:      "0001",
		SpecVersion: ofxgo.OfxVersion220,
	}

	response, err := client.Request(&query)
	if err != nil {
		return nil, errors.Wrap(err, "error requesting OFX transactions")
	}
	if response.Signon.Status.Code != 0 {
		meaning, err := response.Signon.Status.CodeMeaning()
		if err != nil {
			meaning = "UNKNOWN"
		}
		return nil, errors.Errorf("error authenticating with OFX: %d (%s): %s\n", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
	}

	var messages []ofxgo.Message
	switch account.RequestType {
	case ReqTypeBank:
		messages = response.Bank
	case ReqTypeCreditCard:
		messages = response.CreditCard
	default:
		return nil, errors.Errorf("unknown ReqType %q", account.RequestType)
	}
	if len(messages) < 1 {
		return nil, nil
	}
	var transactions []budget.Transaction
	for _, message := range messages {
		if stmt, ok := message.(*ofxgo.StatementResponse); ok {
			txns, err := transactionsFromStatement(stmt)
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, txns...)
		} else if stmt, ok := message.(*ofxgo.CCStatementResponse); ok {
			txns, err := transactionsFromCCStatement(stmt)
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, txns...)
		} else {
			d.Log.WithField("message_type", fmt.Sprintf("%T", message)).
				Warn("Unknown message, ignoring...")
			continue
		}
	}
	budget.TransactionsByDate(transactions)
	return transactions, nil
}

func transactionFromTransaction(tran ofxgo.Transaction) (budget.Transaction, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return budget.Transaction{}, errors.Wrap(err, "Error generating transaction ID: %s")
	}
	amount := big.NewRat(100, 1).Mul(big.NewRat(100, 1), &tran.TrnAmt.Rat)
	if !amount.IsInt() {
		return budget.Transaction{}, errors.Errorf("%s is not specified in whole cents!", tran.TrnAmt.String())
	}
	txn := budget.Transaction{
		ID:              id,
		TransactionType: tran.TrnType.String(),
		DatePosted:      tran.DtPosted.Time,
		Amount:          amount.Num().Int64(),
		FiTID:           string(tran.FiTID),
		ServerTID:       string(tran.SrvrTID),
		CheckNum:        string(tran.CheckNum),
		RefNum:          string(tran.RefNum),
		SIC:             int64(tran.SIC),
		PayeeID:         string(tran.PayeeID),
		Name:            string(tran.Name),
		ExtendedName:    string(tran.ExtdName),
		Memo:            string(tran.Memo),
		CorrectFiTID:    string(tran.CorrectFiTID),
	}
	if tran.Inv401kSource != 0 {
		txn.Inv401kSource = tran.Inv401kSource.String()
	}
	if ok, _ := tran.Currency.Valid(); ok {
		txn.Currency = tran.Currency.CurSym.String()
	}
	if ok, _ := tran.OrigCurrency.Valid(); ok {
		txn.OriginalCurrency = tran.OrigCurrency.CurSym.String()
	}
	if tran.CorrectAction != 0 {
		txn.CorrectAction = tran.CorrectAction.String()
	}
	if tran.DtUser != nil {
		txn.DateUser = tran.DtUser.Time
	}
	if tran.DtAvail != nil {
		txn.DateAvailable = tran.DtAvail.Time
	}
	if tran.Payee != nil {
		txn.PayeeName = string(tran.Payee.Name)
		txn.PayeeAddr1 = string(tran.Payee.Addr1)
		txn.PayeeAddr2 = string(tran.Payee.Addr2)
		txn.PayeeAddr3 = string(tran.Payee.Addr3)
		txn.PayeeCity = string(tran.Payee.City)
		txn.PayeeState = string(tran.Payee.State)
		txn.PayeePostalCode = string(tran.Payee.PostalCode)
		txn.PayeeCountry = string(tran.Payee.Country)
		txn.PayeePhone = string(tran.Payee.Phone)
	}
	if tran.BankAcctTo != nil {
		txn.BankAccountToBankID = string(tran.BankAcctTo.BankID)
		txn.BankAccountToBranchID = string(tran.BankAcctTo.BranchID)
		txn.BankAccountToAccountID = string(tran.BankAcctTo.AcctID)
		txn.BankAccountToAccountType = tran.BankAcctTo.AcctType.String()
		txn.BankAccountToAccountKey = string(tran.BankAcctTo.AcctKey)
	}
	if tran.CCAcctTo != nil {
		txn.CreditCardAccountToAccountID = string(tran.CCAcctTo.AcctID)
		txn.CreditCardAccountToAccountKey = string(tran.CCAcctTo.AcctKey)
	}
	return txn, nil
}

func transactionsFromStatement(stmt *ofxgo.StatementResponse) ([]budget.Transaction, error) {
	results := make([]budget.Transaction, 0, len(stmt.BankTranList.Transactions))
	for _, tran := range stmt.BankTranList.Transactions {
		txn, err := transactionFromTransaction(tran)
		if err != nil {
			return nil, err
		}
		results = append(results, txn)
	}
	return results, nil
}

func transactionsFromCCStatement(stmt *ofxgo.CCStatementResponse) ([]budget.Transaction, error) {
	results := make([]budget.Transaction, 0, len(stmt.BankTranList.Transactions))
	for _, tran := range stmt.BankTranList.Transactions {
		txn, err := transactionFromTransaction(tran)
		if err != nil {
			return nil, err
		}
		results = append(results, txn)
	}
	return results, nil
}
