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
	if g.d.Recurring == nil {
		g.ui.Error("Must have a storer configured for recurring transactions.")
		return 1
	}
	if g.d.Transactions == nil {
		g.ui.Error("Must have a storer configured for transactions.")
		return 1
	}
	groups, err := budget.GroupTransactions(g.ctx, g.d)
	if err != nil {
		yall.FromContext(g.ctx).WithError(err).Error("error grouping transactions")
		return 1
	}
	var recurs []budget.Recurring
	txnChanges := map[string][]string{}
	for _, group := range groups {
		recurs = append(recurs, budget.Recurring{
			ID: group[0].RecurringID,
		})
		for _, txn := range group {
			txnChanges[txn.RecurringID] = append(txnChanges[txn.RecurringID], txn.ID)
		}
	}

	err = g.d.Recurring.CreateRecurrings(g.ctx, recurs)
	if err != nil {
		yall.FromContext(g.ctx).WithError(err).Error("error creating recurring groups")
		return 1
	}
	for recur, txn := range txnChanges {
		recurID := recur
		tf := budget.TransactionFilters{
			IDs: txn,
		}
		change := budget.TransactionChange{
			RecurringID: &recurID,
		}
		yall.FromContext(g.ctx).WithField("group_id", recurID).Info("updating RecurringID on transactions")
		err = g.d.Transactions.UpdateTransactions(g.ctx, tf, change)
		if err != nil {
			yall.FromContext(g.ctx).WithField("group_id", recurID).WithError(err).Error("error updating RecurringID on transactions")
			return 1
		}
	}
	yall.FromContext(g.ctx).Info("updated RecurringID on transactions")
	return 0
}

func (g groupsDetectCommand) Synopsis() string {
	return "Detect transactions that may belong to a recurring group, and update them to be part of it."
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
	if g.d.Recurring == nil {
		yall.FromContext(g.ctx).Error("must have a storer configured for recurring transactions.")
		return 1
	}
	if g.d.Transactions == nil {
		yall.FromContext(g.ctx).Error("must have a storer configured for transactions.")
		return 1
	}
	if len(args) < 1 {
		yall.FromContext(g.ctx).Error("must specify a recurring ID")
		return 1
	}
	log := yall.FromContext(g.ctx).WithField("recurring_id", args[0])
	log = log.WithField("transactions_storer", fmt.Sprintf("%T", g.d.Transactions))
	log = log.WithField("groups_storer", fmt.Sprintf("%T", g.d.Recurring))
	transactions, err := g.d.Transactions.ListTransactions(g.ctx, budget.TransactionFilters{
		RecurringID: &args[0],
	})
	if err != nil {
		log.WithError(err).Error("error retrieving transactions")
		return 1
	}
	if len(transactions) < 1 {
		log.Error("no transactions for that recurring ID")
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
	log.WithField("total_days", totalDays).WithField("populated_days", len(days)).Info("grouped transactions by day")

	transactionsPerDay := make([]float64, totalDays+1)
	centsPerDay := make([]float64, totalDays+1)
	for day, txns := range days {
		log.WithField("day", day).Debug("processing transactions per day")
		transactionsPerDay[day] = float64(len(txns))
		for _, txn := range txns {
			centsPerDay[day] += float64(txn.Amount)
		}
	}
	txnDensities, txnFreqs := spectral.Pwelch(transactionsPerDay, 1.0, &spectral.PwelchOptions{})
	txnPSD := zipDensitiesAndFrequencies(txnDensities, txnFreqs)
	sort.Slice(txnPSD, func(i, j int) bool { return txnPSD[i].density > txnPSD[j].density })
	for _, psd := range txnPSD {
		log.WithField("density", psd.density).WithField("frequency", psd.frequency).Debug("calculated PSD for transactions")
	}

	centDensities, centFreqs := spectral.Pwelch(centsPerDay, 1.0, &spectral.PwelchOptions{})
	centPSD := zipDensitiesAndFrequencies(centDensities, centFreqs)
	sort.Slice(centPSD, func(i, j int) bool { return centPSD[i].density > centPSD[j].density })
	for _, psd := range centPSD {
		log.WithField("density", psd.density).WithField("frequency", psd.frequency).Debug("calculated PSD for amount")
	}
	log.WithField("days_in_period", math.Round(1.0/txnPSD[0].frequency)).Info("calculated period")
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
	return "Forecast transactions that may belong to a recurring group, and update them to be part of it."
}
