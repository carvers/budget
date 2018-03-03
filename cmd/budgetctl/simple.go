package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carvers/budget"
	"github.com/carvers/budget/simple"
	"github.com/mitchellh/cli"
	yall "yall.in"
)

func simpleImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return simpleImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type simpleImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (s simpleImportCommand) Help() string {
	return "help text"
}

func (s simpleImportCommand) Run(args []string) int {
	log := yall.FromContext(s.ctx)
	if len(args) < 1 {
		log.Error("no filename specified")
		return 1
	}
	log = log.WithField("file", os.Args[0])

	if s.d.Transactions == nil {
		log.Error("must specify transaction store")
		return 1
	}
	if s.d.Accounts == nil {
		log.Error("must specify account store")
		return 1
	}

	file, err := os.Open(os.Args[0])
	if err != nil {
		log.WithError(err).Error("error opening file")
		return 1
	}

	transactions, err := simple.FromReader(s.ctx, file)
	if err != nil {
		log.WithError(err).Error("error parsing Simple transactions")
		return 1
	}
	if len(transactions) < 1 {
		log.Warn("no transactions")
		os.Exit(1)
	}
	accountID := transactions[0].AccountID
	log = log.WithField("account", accountID).WithField("accounts_storer", fmt.Sprintf("%T", s.d.Accounts))
	log.Info("persisting account")

	err = s.d.Accounts.CreateAccount(s.ctx, budget.Account{
		ID:          accountID,
		AccountType: "CHECKING",
		Sync:        false,
	})
	log.Info("persisted account")
	log = log.WithField("transactions_storer", fmt.Sprintf("%T", s.d.Transactions))
	for i := 0; i < (len(transactions)/50)+1; i++ {
		if len(transactions) <= i*50 {
			continue
		}
		txns := transactions[i*50:]
		if len(txns) > 50 {
			txns = txns[:50]
		}
		log := log.WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions)))
		log.Info("persisting transactions")

		// upsert those transactions into postgres
		err = s.d.Transactions.ImportTransactions(s.ctx, txns)
		if err != nil {
			log.WithError(err).Error("error saving transactions")
			return 1
		}
		log.Info("persisted transactions")
	}
	return 0
}

func (s simpleImportCommand) Synopsis() string {
	return "Load in downloaded transactions from Simple in JSON format."
}
