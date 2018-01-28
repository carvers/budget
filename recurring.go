package budget

import "impractical.co/pqarrays"

type Recurring struct {
	ID         string
	Name       string
	AccountIDs pqarrays.StringArray
}

func (r Recurring) GetSQLTableName() string {
	return "recurring_groups"
}

type RecurringChange struct {
	Name       *string
	AccountIDs []string
}

func (c RecurringChange) IsEmpty() bool {
	if c.Name != nil {
		return false
	}
	if c.AccountIDs != nil {
		return false
	}
	return true
}
