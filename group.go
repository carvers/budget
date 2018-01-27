package budget

import (
	"github.com/carvers/budget/similar"
)

func GroupTransactions(d Dependencies) ([][]Transaction, error) {
	transactions, err := d.Transactions.ListTransactions(TransactionFilters{})
	if err != nil {
		return nil, err
	}
	uniqueTransactions := map[string][]Transaction{}
	for _, txn := range transactions {
		name := similar.Sanitize(txn.Name)
		uniqueTransactions[name] = append(uniqueTransactions[name], txn)
	}
	var groups [][]Transaction
	for name, txns := range uniqueTransactions {
		bestMatch := -1
		bestJac := 0.0
		for pos, group := range groups {
			for _, member := range group {
				memberName := similar.Sanitize(member.Name)
				jac := similar.Jaccard(similar.Shingle(memberName), similar.Shingle(name))
				if jac > bestJac {
					bestMatch = pos
				}
				if jac > .9 {
					break
				}
			}
			if bestJac > .9 {
				break
			}
		}
		if bestJac >= 0.4 {
			groups[bestMatch] = append(groups[bestMatch], txns...)
		} else {
			groups = append(groups, txns)
		}
	}
	var results [][]Transaction
	for _, group := range groups {
		if len(group) <= 1 {
			continue
		}
		results = append(results, group)
	}
	return results, nil
}
