package budget

import (
	"context"

	"github.com/carvers/budget/similar"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"impractical.co/pqarrays"
)

type Group struct {
	ID         string
	Name       string
	AccountIDs pqarrays.StringArray
	TrendID    string
	Finished   bool
}

func (g Group) GetSQLTableName() string {
	return "transaction_groups"
}

type GroupChange struct {
	Name       *string
	AccountIDs []string
	Finished   *bool
}

func (c GroupChange) IsEmpty() bool {
	if c.Name != nil {
		return false
	}
	if c.AccountIDs != nil {
		return false
	}
	return true
}

func GroupTransactions(ctx context.Context, d Dependencies) ([][]Transaction, error) {
	// retrieve a list of all our transactions
	transactions, err := d.Transactions.ListTransactions(ctx, TransactionFilters{})
	if err != nil {
		return nil, errors.Wrap(err, "error listing transactions")
	}

	// retrieve a list of the transaction groups that have similar features
	groups, err := d.Groups.ListGroups(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error listing groups")
	}
	// ID-index the groups for easy reference
	groupMap := make(map[string]Group, len(groups))
	for _, r := range groups {
		groupMap[r.ID] = r
	}

	// create a map of BudgetID -> slice of transactions so we can sort out
	// groups of GroupID
	txnGroups := map[string][]Transaction{}

	// if they have the same name, we don't need to evaluate the name
	// multiple times, so just group them
	uniqueTransactions := map[string][]Transaction{}
	for _, txn := range transactions {
		// if we already have a group ID, we don't need to evaluate
		// it at all.
		if txn.GroupID != "" {
			txnGroups[txn.GroupID] = append(txnGroups[txn.GroupID], txn)
			continue
		}
		// otherwise, add it to the list of things to examine
		// we need to prefix the name with the account ID, so we
		// are considering only duplicates within an account duplicates
		// because some groups may have account restrictions for
		// which accounts they match.
		name := txn.AccountID + ":" + similar.Sanitize(txn.Name)
		uniqueTransactions[name] = append(uniqueTransactions[name], txn)
	}

	for _, txns := range uniqueTransactions {
		// get the name for this group
		name := similar.Sanitize(txns[0].Name)

		// keep track of what the best score is, and which groupID
		// scored it
		var bestJac float64
		var groupID string

		// Loop through all the groups of GroupIDs that we already
		// have, and check if the name fits with that group.
		for id, group := range txnGroups {
			// if that group has account restrictions, and we're not
			// in one of those accounts, don't consider that group.
			if g, ok := groupMap[id]; ok && len(g.AccountIDs) > 0 {
				var inAccounts bool
				for _, account := range g.AccountIDs {
					if account == txns[0].AccountID {
						inAccounts = true
						break
					}
				}
				if !inAccounts {
					continue
				}
			}

			// compare our name with each transaction already in the
			// group
			for _, member := range group {
				memberName := similar.Sanitize(member.Name)
				// if the name is an exact match, just skip all this
				if memberName == name {
					groupID = id
					bestJac = 1.0
					break
				}

				jac := similar.Jaccard(similar.Shingle(memberName), similar.Shingle(name))
				// if our best score so far is less than the score
				// for the transaction we just matched against, the
				// name belongs with that transaction.
				if jac > bestJac {
					groupID = id
				}
				// if our jaccard coefficient is over 90%, the
				// odds of us getting a better match are pretty
				// low, so let's just accept that as the answer.
				// This saves us from having to compare every
				// single transaction against every other one.
				if jac > .9 {
					break
				}
			}
			// to accept it as an answer, we need to break out of the
			// loop that is comparing against all the transactions in
			// the group, but also break out of the loop searching all
			// the transaction groups.
			if bestJac > .9 {
				break
			}
		}
		// if our score is less than .4, we don't consider it a match and
		// need to create a new group for this transaction name
		if bestJac < 0.4 {
			groupID, err = uuid.GenerateUUID()
			if err != nil {
				return nil, errors.Wrap(err, "error generating UUID")
			}
		}
		// make sure every transaction in the group has the group ID set
		for pos, txn := range txns {
			txn.GroupID = groupID
			txns[pos] = txn
		}
		// add our transactions to the group, and carry on to the next
		// transaction name.
		txnGroups[groupID] = append(txnGroups[groupID], txns...)
	}
	// turn our map into a slice of slices to return
	var results [][]Transaction
	for _, group := range txnGroups {
		// if we have a group of 1, that's not _really_ a group, so we
		// just assume that's not a group.
		if len(group) <= 1 {
			continue
		}
		results = append(results, group)
	}
	return results, nil
}
