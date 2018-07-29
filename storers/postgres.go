package storers

import (
	"context"
	"database/sql"

	"darlinggo.co/pan"
	yall "yall.in"

	"github.com/carvers/budget"
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) postgres {
	return postgres{db: db}
}

func importTransactionsSQL(t []budget.Transaction) *pan.Query {
	pannable := make([]pan.SQLTableNamer, 0, len(t))
	for _, txn := range t {
		pannable = append(pannable, txn)
	}
	q := pan.Insert(pannable...)
	q.Expression("ON CONFLICT (")
	q.Expression(pan.Column(t[0], "AccountID") + ",")
	q.Expression(pan.Column(t[0], "FiTID") + " )")
	q.Expression("DO NOTHING")
	return q.Flush(" ")
}

func (p postgres) ImportTransactions(ctx context.Context, t []budget.Transaction) error {
	query := importTransactionsSQL(t)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("Importing transactions")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func listTransactionsSQL(f budget.TransactionFilters) *pan.Query {
	var t budget.Transaction
	q := pan.New("SELECT " + pan.Columns(t).String() + " FROM " + pan.Table(t))
	addTransactionFiltersToQuery(q, f)
	q.OrderBy("date_posted DESC, account_id, amount DESC")
	return q.Flush(" ")
}

func (p postgres) ListTransactions(ctx context.Context, f budget.TransactionFilters) ([]budget.Transaction, error) {
	query := listTransactionsSQL(f)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return nil, err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("Listing transactions")
	rows, err := p.db.Query(queryStr, query.Args()...)
	if err != nil {
		return nil, err
	}
	var transactions []budget.Transaction
	for rows.Next() {
		var t budget.Transaction
		err = pan.Unmarshal(rows, &t)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return transactions, err
	}
	return transactions, nil
}

func updateTransactionsSQL(tf budget.TransactionFilters, change budget.TransactionChange) *pan.Query {
	var t budget.Transaction
	q := pan.New("UPDATE " + pan.Table(t) + " SET")
	if change.GroupID != nil {
		q.Comparison(t, "GroupID", "=", *change.GroupID)
	}
	q.Flush(", ")
	addTransactionFiltersToQuery(q, tf)
	return q
}

func (p postgres) UpdateTransactions(ctx context.Context, tf budget.TransactionFilters, change budget.TransactionChange) error {
	if change.IsEmpty() {
		return nil
	}
	query := updateTransactionsSQL(tf, change)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("updating transactions")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func balanceSQL(account string) *pan.Query {
	var t budget.Transaction
	q := pan.New("SELECT SUM(" + pan.Column(t, "Amount") + ") FROM " + pan.Table(t))
	q.Where()
	q.Comparison(t, "AccountID", "=", account)
	return q.Flush(" ")
}

func (p postgres) Balance(ctx context.Context, accountID string) (int64, error) {
	var bal int64
	query := balanceSQL(accountID)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return 0, err
	}
	yall.FromContext(ctx).WithField("query", queryStr).WithField("account", accountID).Debug("retrieving balance")
	rows, err := p.db.Query(queryStr, query.Args()...)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = pan.Unmarshal(rows, &bal)
		if err != nil {
			return bal, err
		}
	}
	if err = rows.Err(); err != nil {
		return bal, err
	}
	return bal, nil
}

func createAccountSQL(account budget.Account) *pan.Query {
	return pan.Insert(account)
}

func (p postgres) CreateAccount(ctx context.Context, account budget.Account) error {
	query := createAccountSQL(account)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("Creating account")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func getAccountSQL(id string) *pan.Query {
	var account budget.Account
	q := pan.New("SELECT " + pan.Columns(account).String() + " FROM " + pan.Table(account))
	q.Where()
	q.Comparison(account, "ID", "=", id)
	return q.Flush(" ")
}

func (p postgres) GetAccount(ctx context.Context, id string) (budget.Account, error) {
	query := getAccountSQL(id)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return budget.Account{}, err
	}
	rows, err := p.db.Query(queryStr, query.Args()...)
	if err != nil {
		return budget.Account{}, err
	}
	var account budget.Account
	for rows.Next() {
		err = pan.Unmarshal(rows, &account)
		if err != nil {
			return account, err
		}
	}
	if err = rows.Err(); err != nil {
		return account, err
	}
	if account.ID == "" {
		return account, budget.ErrAccountNotFound
	}
	return account, nil
}

func updateAccountSQL(id string, change budget.AccountChange) *pan.Query {
	var account budget.Account
	q := pan.New("UPDATE " + pan.Table(account) + " SET ")
	if change.Name != nil {
		q.Comparison(account, "Name", "=", *change.Name)
	}
	if change.AccountType != nil {
		q.Comparison(account, "AccountType", "=", *change.AccountType)
	}
	if change.RequestType != nil {
		q.Comparison(account, "RequestType", "=", *change.RequestType)
	}
	if change.BankORG != nil {
		q.Comparison(account, "BankORG", "=", *change.BankORG)
	}
	if change.BankFID != nil {
		q.Comparison(account, "BankFID", "=", *change.BankFID)
	}
	if change.BankURL != nil {
		q.Comparison(account, "BankURL", "=", *change.BankURL)
	}
	q.Flush(", ")
	q.Where()
	q.Comparison(account, "ID", "=", id)
	return q.Flush(" ")
}

func (p postgres) UpdateAccount(ctx context.Context, id string, change budget.AccountChange) error {
	if change.IsEmpty() {
		return nil
	}
	query := updateAccountSQL(id, change)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).WithField("id", id).Debug("updating account")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func deleteAccountSQL(id string) *pan.Query {
	var account budget.Account
	q := pan.New("DELETE FROM " + pan.Table(account))
	q.Where()
	q.Comparison(account, "ID", "=", id)
	return q.Flush(" ")
}

func (p postgres) DeleteAccount(ctx context.Context, id string) error {
	query := deleteAccountSQL(id)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).WithField("id", id).Debug("deleting account")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func listAccountsSQL() *pan.Query {
	var account budget.Account
	q := pan.New("SELECT " + pan.Columns(account).String() + " FROM " + pan.Table(account))
	q.OrderByDesc(pan.Column(account, "Name"))
	return q.Flush(" ")
}

func (p postgres) ListAccounts(ctx context.Context) ([]budget.Account, error) {
	query := listAccountsSQL()
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return nil, err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("listing account")
	rows, err := p.db.Query(queryStr, query.Args()...)
	if err != nil {
		return nil, err
	}
	var accts []budget.Account
	for rows.Next() {
		var account budget.Account
		err = pan.Unmarshal(rows, &account)
		if err != nil {
			return accts, err
		}
		accts = append(accts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accts, nil
}

func createGroupsSQL(groups []budget.Group) *pan.Query {
	pannable := make([]pan.SQLTableNamer, 0, len(groups))
	for _, g := range groups {
		pannable = append(pannable, g)
	}
	q := pan.Insert(pannable...)
	q.Expression("ON CONFLICT ON CONSTRAINT " + pan.Table(groups[0]) + "_pkey")
	q.Expression("DO NOTHING")
	return q.Flush(" ")
}

func (p postgres) CreateGroups(ctx context.Context, groups []budget.Group) error {
	query := createGroupsSQL(groups)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).WithField("num_groups", len(groups)).
		Debug("creating groups")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func listGroupsSQL() *pan.Query {
	var group budget.Group
	q := pan.New("SELECT " + pan.Columns(group).String() + " FROM " + pan.Table(group))
	q.OrderByDesc(pan.Column(group, "ID"))
	return q.Flush(" ")
}

func (p postgres) ListGroups(ctx context.Context) ([]budget.Group, error) {
	query := listGroupsSQL()
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return nil, err
	}
	yall.FromContext(ctx).WithField("query", queryStr).Debug("listing groups")
	rows, err := p.db.Query(queryStr, query.Args()...)
	if err != nil {
		return nil, err
	}
	var groups []budget.Group
	for rows.Next() {
		var group budget.Group
		err = pan.Unmarshal(rows, &group)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

func updateGroupSQL(id string, change budget.GroupChange) *pan.Query {
	// TODO(paddy): write SQL for updating a group
	return nil
}

func (p postgres) UpdateGroup(ctx context.Context, id string, change budget.GroupChange) error {
	if change.IsEmpty() {
		return nil
	}
	query := updateGroupSQL(id, change)
	queryStr, err := query.PostgreSQLString()
	if err != nil {
		return err
	}
	yall.FromContext(ctx).WithField("query", queryStr).WithField("id", id).Debug("updating group")
	_, err = p.db.Exec(queryStr, query.Args()...)
	return err
}

func addTransactionFiltersToQuery(q *pan.Query, f budget.TransactionFilters) {
	var t budget.Transaction
	if len(f.IDs) > 0 {
		q.Where()
		iface := make([]interface{}, 0, len(f.IDs))
		for _, id := range f.IDs {
			iface = append(iface, id)
		}
		q.In(t, "ID", iface...)
	}
	if f.AccountIDs != nil {
		if len(f.AccountIDs) > 0 {
			q.Where()
			iface := make([]interface{}, 0, len(f.AccountIDs))
			for _, id := range f.AccountIDs {
				iface = append(iface, id)
			}
			q.In(t, "AccountID", iface...)
		} else {
			// specifically set an empty array, meaning not set
			q.Where()
			q.Comparison(t, "AccountID", "=", "")
		}
	}
	if f.TransactionTypes != nil {
		if len(f.TransactionTypes) > 0 {
			q.Where()
			iface := make([]interface{}, 0, len(f.TransactionTypes))
			for _, typ := range f.TransactionTypes {
				iface = append(iface, typ)
			}
			q.In(t, "TransactionType", iface...)
		} else {
			// specifically set an empty array, meaning not set
			q.Where()
			q.Comparison(t, "TransactionType", "=", "")
		}
	}
	if f.DatePostedBefore != nil {
		q.Where()
		q.Comparison(t, "DatePosted", "<", *f.DatePostedBefore)
	}
	if f.DatePostedAfter != nil {
		q.Where()
		q.Comparison(t, "DatePosted", ">", *f.DatePostedAfter)
	}
	if f.Amount != nil {
		q.Where()
		q.Comparison(t, "Amount", "=", *f.Amount)
	}
	if f.AmountGreaterThan != nil {
		q.Where()
		q.Comparison(t, "Amount", ">", *f.AmountGreaterThan)
	}
	if f.AmountLessThan != nil {
		q.Where()
		q.Comparison(t, "Amount", "<", *f.AmountLessThan)
	}
	if f.CheckNum != nil {
		q.Where()
		q.Comparison(t, "CheckNum", "=", *f.CheckNum)
	}
	if f.RefNum != nil {
		q.Where()
		q.Comparison(t, "RefNum", "=", *f.RefNum)
	}
	if f.Name != nil {
		q.Where()
		q.Comparison(t, "Name", "=", *f.Name)
	}
	if f.GroupID != nil {
		q.Where()
		q.Comparison(t, "GroupID", "=", *f.GroupID)
	}
	q.Flush(" AND ")
}
