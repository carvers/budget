package apiv1

import "github.com/carvers/budget"

type Group struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	AccountIDs []string `json:"account_ids"`
	TrendID    string   `json:"trend_id"`
	Finished   bool     `json:"finished"`
}

func groupFromCore(group budget.Group) Group {
	return Group{
		ID:         group.ID,
		Name:       group.Name,
		AccountIDs: []string(group.AccountIDs),
		TrendID:    group.TrendID,
		Finished:   group.Finished,
	}
}

func groupsFromCore(groups []budget.Group) []Group {
	res := make([]Group, 0, len(groups))
	for _, group := range groups {
		res = append(res, groupFromCore(group))
	}
	return res
}
