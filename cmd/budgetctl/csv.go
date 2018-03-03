package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carvers/budget"
	"github.com/carvers/budget/csv"
	"github.com/mitchellh/cli"
	yall "yall.in"
)

func csvImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return csvImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type csvImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (c csvImportCommand) Help() string {
	return "help text"
}

func (c csvImportCommand) Run(args []string) int {
	if c.d.Transactions == nil {
		yall.FromContext(c.ctx).Error("must have a transaction storer configured")
		return 1
	}
	if len(args) < 2 {
		yall.FromContext(c.ctx).Error("The csv import command expects two arguments")
		return 1
	}
	file, err := os.Open(args[1])
	if err != nil {
		yall.FromContext(c.ctx).WithField("file", args[1]).WithError(err).Error("error opening file")
		return 1
	}
	defer file.Close()
	transactions, err := csv.FromReader(file, args[0])
	if err != nil {
		yall.FromContext(c.ctx).WithField("file", args[1]).WithError(err).Error("error reading transactions from file")
		return 1
	}
	if len(transactions) < 1 {
		yall.FromContext(c.ctx).Warn("No transactions in CSV")
		os.Exit(2)
	}
	for i := 0; i < (len(transactions)/50)+1; i++ {
		if len(transactions) <= i*50 {
			continue
		}
		txns := transactions[i*50:]
		if len(txns) > 50 {
			txns = txns[:50]
		}
		yall.FromContext(c.ctx).WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", c.d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisting transactions")

		// upsert those transactions into postgres
		err = c.d.Transactions.ImportTransactions(c.ctx, txns)
		if err != nil {
			yall.FromContext(c.ctx).WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", c.d.Transactions)).
				WithError(err).Error("error saving transactions")
			return 1
		}
		yall.FromContext(c.ctx).WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", c.d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisted transactions")
	}
	return 0
}

func (c csvImportCommand) Synopsis() string {
	return "Load in downloaded transactions in CSV format."
}
