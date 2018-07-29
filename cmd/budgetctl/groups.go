package main

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
	"github.com/mjibson/go-dsp/spectral"
	yall "yall.in"
)

func groupsDetectCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return groupsDetectCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type groupsDetectCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (g groupsDetectCommand) Help() string {
	return "help text"
}

func (g groupsDetectCommand) Run(args []string) int {
	if g.d.Groups == nil {
		g.ui.Error("Must have a storer configured for transaction groups.")
		return 1
	}
	if g.d.Transactions == nil {
		g.ui.Error("Must have a storer configured for transactions.")
		return 1
	}
	txnGroups, err := budget.GroupTransactions(g.ctx, g.d)
	if err != nil {
		yall.FromContext(g.ctx).WithError(err).Error("error grouping transactions")
		return 1
	}
	var groups []budget.Group
	txnChanges := map[string][]string{}
	for _, group := range txnGroups {
		groups = append(groups, budget.Group{
			ID: group[0].GroupID,
		})
		for _, txn := range group {
			txnChanges[txn.GroupID] = append(txnChanges[txn.GroupID], txn.ID)
		}
	}

	err = g.d.Groups.CreateGroups(g.ctx, groups)
	if err != nil {
		yall.FromContext(g.ctx).WithError(err).Error("error creating groups")
		return 1
	}
	for group, txn := range txnChanges {
		groupID := group
		tf := budget.TransactionFilters{
			IDs: txn,
		}
		change := budget.TransactionChange{
			GroupID: &groupID,
		}
		yall.FromContext(g.ctx).WithField("group_id", groupID).Info("updating GroupID on transactions")
		err = g.d.Transactions.UpdateTransactions(g.ctx, tf, change)
		if err != nil {
			yall.FromContext(g.ctx).WithField("group_id", groupID).WithError(err).Error("error updating GroupID on transactions")
			return 1
		}
	}
	yall.FromContext(g.ctx).Info("updated GroupID on transactions")
	return 0
}

func (g groupsDetectCommand) Synopsis() string {
	return "Detect transactions that may belong to a group, and update them to be part of it."
}

func groupsPeriodsCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return groupsPeriodsCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type groupsPeriodsCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (g groupsPeriodsCommand) Help() string {
	return "help text"
}

func (g groupsPeriodsCommand) Run(args []string) int {
	if g.d.Groups == nil {
		yall.FromContext(g.ctx).Error("must have a storer configured for transaction groups.")
		return 1
	}
	if g.d.Transactions == nil {
		yall.FromContext(g.ctx).Error("must have a storer configured for transactions.")
		return 1
	}
	if len(args) < 1 {
		yall.FromContext(g.ctx).Error("must specify a group ID")
		return 1
	}
	log := yall.FromContext(g.ctx).WithField("group_id", args[0])
	log = log.WithField("transactions_storer", fmt.Sprintf("%T", g.d.Transactions))
	log = log.WithField("groups_storer", fmt.Sprintf("%T", g.d.Groups))
	transactions, err := g.d.Transactions.ListTransactions(g.ctx, budget.TransactionFilters{
		GroupID: &args[0],
	})
	if err != nil {
		log.WithError(err).Error("error retrieving transactions")
		return 1
	}
	if len(transactions) < 1 {
		log.Error("no transactions for that group ID")
		return 1
	}

	sort.Slice(transactions, func(i, j int) bool { return transactions[i].DatePosted.Before(transactions[j].DatePosted) })

	days := map[int64][]budget.Transaction{}
	var totalDays int64
	for _, txn := range transactions {
		offset := int64(txn.DatePosted.Sub(transactions[0].DatePosted).Truncate(time.Hour) / time.Hour / 24)
		days[offset] = append(days[offset], txn)
		if offset > totalDays {
			totalDays = offset
		}
	}
	log.WithField("total_days", totalDays).WithField("populated_days", len(days)).Debug("grouped transactions by day")

	transactionsPerDay := make([]float64, totalDays+1)
	for day, txns := range days {
		log.WithField("day", day).Debug("processing transactions per day")
		transactionsPerDay[day] = float64(len(txns))
	}
	txnDensities, txnFreqs := spectral.Pwelch(transactionsPerDay, 1.0, &spectral.PwelchOptions{})
	txnPSD := zipDensitiesAndFrequencies(txnDensities, txnFreqs)
	sort.Slice(txnPSD, func(i, j int) bool { return txnPSD[i].density > txnPSD[j].density })
	for _, psd := range txnPSD {
		log.WithField("density", psd.density).WithField("frequency", psd.frequency).Debug("calculated PSD for transactions")
	}

	var freq float64
	for _, psd := range txnPSD {
		freq = psd.frequency
		if freq != 0 {
			break
		}
	}
	log.WithField("days_in_period", math.Round(1.0/freq)).Info("calculated period")
	return 0
}

type psd struct {
	density   float64
	frequency float64
}

func zipDensitiesAndFrequencies(d, f []float64) []psd {
	res := make([]psd, 0, len(d))
	for pos, den := range d {
		res = append(res, psd{
			density:   den,
			frequency: f[pos],
		})
	}
	return res
}

func (g groupsPeriodsCommand) Synopsis() string {
	return "Forecast transactions that may belong to a group, and update them to be part of it."
}
