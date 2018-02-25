package budget

import (
	"context"

	"github.com/carvers/budget/similar"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
)

func GroupTransactions(ctx context.Context, d Dependencies) ([][]Transaction, error) {
	// retrieve a list of all our transactions
	transactions, err := d.Transactions.ListTransactions(ctx, TransactionFilters{})
	if err != nil {
		return nil, errors.Wrap(err, "error listing transactions")
	}

	// retrieve a list of the recurring groups
	recurring, err := d.Recurring.ListRecurrings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error listing recurring groups")
	}
	// ID-index the recurring groups for easy reference
	recurMap := make(map[string]Recurring, len(recurring))
	for _, r := range recurring {
		recurMap[r.ID] = r
	}

	// create a map of BudgetID -> slice of transactions so we can sort out
	// groups of RecurringID
	groups := map[string][]Transaction{}

	// if they have the same name, we don't need to evaluate the name
	// multiple times, so just group them
	uniqueTransactions := map[string][]Transaction{}
	for _, txn := range transactions {
		// if we already have a recurring ID, we don't need to evaluate
		// it at all.
		if txn.RecurringID != "" {
			groups[txn.RecurringID] = append(groups[txn.RecurringID], txn)
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

		// keep track of what the best score is, and which recurringID
		// scored it
		var bestJac float64
		var recurringID string

		// Loop through all the groups of RecurringIDs that we already
		// have, and check if the name fits with that group.
		for id, group := range groups {
			// if that group has account restrictions, and we're not
			// in one of those accounts, don't consider that group.
			if g, ok := recurMap[id]; ok && len(g.AccountIDs) > 0 {
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
					recurringID = id
					bestJac = 1.0
					break
				}

				jac := similar.Jaccard(similar.Shingle(memberName), similar.Shingle(name))
				// if our best score so far is less than the score
				// for the transaction we just matched against, the
				// name belongs with that transaction.
				if jac > bestJac {
					recurringID = id
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
			// the recurring group, but also break out of the loop
			// searching all the recurring groups.
			if bestJac > .9 {
				break
			}
		}
		// if our score is less than .4, we don't consider it a match and
		// need to create a new group for this transaction name
		if bestJac < 0.4 {
			recurringID, err = uuid.GenerateUUID()
			if err != nil {
				return nil, errors.Wrap(err, "error generating UUID")
			}
		}
		// make sure every transaction in the group has the recurring ID set
		for pos, txn := range txns {
			txn.RecurringID = recurringID
			txns[pos] = txn
		}
		// add our transactions to the group, and carry on to the next
		// transaction name.
		groups[recurringID] = append(groups[recurringID], txns...)
	}
	// turn our map into a slice of slices to return
	var results [][]Transaction
	for _, group := range groups {
		// if we have a group of 1, that's not _really_ a group, so we
		// just assume that's not a recurring transaction.
		if len(group) <= 1 {
			continue
		}
		results = append(results, group)
	}
	return results, nil
}
