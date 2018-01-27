package ofx

/*
func ImportFromFile(path, institution string) ([]Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	resp, err := ofxgo.ParseResponse(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	var transactions []Transaction
	for _, msg := range resp.CreditCard {
		switch msg.(type) {
		case *ofxgo.CCStatementResponse:
			stmt := msg.(*ofxgo.CCStatementResponse)
			for _, t := range stmt.BankTranList.Transactions {
				txn := TransactionFromOFX(t)
				txn.CreditCardAccountFromAccountID = stmt.CCAcctFrom.AcctID.String()
				txn.CreditCardAccountFromAccountKey = stmt.CCAcctFrom.AcctKey.String()
				txn.AccountID = txn.CreditCardAccountFromAccountID
				txn.Institution = institution
				transactions = append(transactions, txn)
			}
		default:
			return nil, fmt.Errorf("unknown message type %T", msg)
		}
	}
	for _, msg := range resp.Bank {
		switch msg.(type) {
		case *ofxgo.StatementResponse:
			stmt := msg.(*ofxgo.StatementResponse)
			for _, t := range stmt.BankTranList.Transactions {
				txn := TransactionFromOFX(t)
				transactions = append(transactions, txn)
			}
		default:
			return nil, fmt.Errorf("unknown message type %T", msg)
		}
	}
	return transactions, nil
}
*/
