package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carvers/budget"
	"github.com/carvers/budget/ofx"
	"github.com/mitchellh/cli"
	yall "yall.in"
)

func ofxFetchCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return ofxFetchCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type ofxFetchCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (o ofxFetchCommand) Help() string {
	return "help text"
}

func (o ofxFetchCommand) Run(args []string) int {
	if o.d.AccountsSensitive == nil {
		yall.FromContext(o.ctx).Error("must have a sensitive data store configured to fetch transactions")
		return 1
	}
	if o.d.Transactions == nil {
		yall.FromContext(o.ctx).Error("must have a transaction store configured")
		return 1
	}
	yall.FromContext(o.ctx).Info("retrieving accounts")
	accounts, err := o.d.Accounts.ListAccounts(o.ctx)
	if err != nil {
		yall.FromContext(o.ctx).WithError(err).Error("error listing accounts")
		return 1
	}
	for _, account := range accounts {
		if !account.Sync {
			yall.FromContext(o.ctx).WithField("account", account.ID).WithField("name", account.Name).Info("not configured to sync, ignoring")
			continue
		}
		yall.FromContext(o.ctx).WithField("account", account.ID).WithField("name", account.Name).WithField("storer", fmt.Sprintf("%T", o.d.AccountsSensitive)).Info("fetching account's sensitive details")
		asd, err := o.d.AccountsSensitive.GetAccountSensitiveDetails(o.ctx, account.ID)
		if err != nil {
			yall.FromContext(o.ctx).WithField("account", account.ID).WithError(err).Error("error retrieving sensitive details")
			continue
		}

		yall.FromContext(o.ctx).WithField("account", account.ID).WithField("name", account.Name).Info("fetching account transactions")
		transactions, err := ofx.FetchTransactions(o.ctx, o.d, account, asd, "3b56383c-2cb6-4c01-8ce0-4951caaf4fa5")
		if err != nil {
			yall.FromContext(o.ctx).WithField("account", account.ID).WithError(err).Error("error fetching transactions")
			continue
		}

		for pos, txn := range transactions {
			txn.AccountID = account.ID
			transactions[pos] = txn
		}
		yall.FromContext(o.ctx).WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", o.d.Transactions)).WithField("num_transactions", len(transactions)).Info("persisting transactions")

		err = o.d.Transactions.ImportTransactions(o.ctx, transactions)
		if err != nil {
			yall.FromContext(o.ctx).WithError(err).Error("error saving transactions")
			return 1
		}
		yall.FromContext(o.ctx).WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", o.d.Transactions)).WithField("num_transactions", len(transactions)).Info("persisted transactions")
	}
	return 0
}

func (o ofxFetchCommand) Synopsis() string {
	return "Load in recent transactions from your accounts using OFX."
}

func ofxImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return ofxImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type ofxImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (o ofxImportCommand) Help() string {
	return "help text"
}

func (o ofxImportCommand) Run(args []string) int {
	if o.d.Transactions == nil {
		yall.FromContext(o.ctx).Error("must have a transaction storer configured")
		return 1
	}
	if o.d.Accounts == nil {
		yall.FromContext(o.ctx).Error("must have a account storer configured")
		return 1
	}
	if o.d.AccountsSensitive == nil {
		yall.FromContext(o.ctx).Error("must have a storer for sensitive account data configured")
		return 1
	}
	if len(args) < 1 {
		yall.FromContext(o.ctx).Error("must specify a file to import")
		return 1
	}
	log := yall.FromContext(o.ctx).WithField("file", args[0])
	file, err := os.Open(args[0])
	if err != nil {
		log.WithError(err).Error("error opening file")
		return 1
	}
	asd, transactions, err := ofx.FromReader(o.ctx, file)
	if err != nil {
		log.WithError(err).Error("error parsing OFX transactions")
		return 1
	}
	if len(transactions) < 1 {
		log.Warn("no transactions")
		return 2
	}
	log.WithField("storer", fmt.Sprintf("%T", o.d.Accounts)).Info("fetching accounts")
	accounts, err := o.d.Accounts.ListAccounts(o.ctx)
	if err != nil {
		log.WithError(err).Error("error listing accounts")
		return 1
	}

	var accountID string
	for _, account := range accounts {
		log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", o.d.AccountsSensitive)).Info("fetching accounts sensitive details")
		a, err := o.d.AccountsSensitive.GetAccountSensitiveDetails(o.ctx, account.ID)
		if err != nil {
			log.WithField("account", account.ID).WithError(err).Error("error retrieving sensitive details")
			continue
		}
		if a.AccountID != asd.AccountID {
			continue
		}
		if a.BankID != asd.BankID {
			continue
		}
		accountID = account.ID
	}
	if accountID == "" {
		log.Error("unable to find account")
		return 1
	}
	log = log.WithField("account", accountID)
	for pos, txn := range transactions {
		txn.AccountID = accountID
		transactions[pos] = txn
	}

	for i := 0; i < (len(transactions)/50)+1; i++ {
		if len(transactions) <= i*50 {
			continue
		}
		txns := transactions[i*50:]
		if len(txns) > 50 {
			txns = txns[:50]
		}
		log := log.WithField("storer", fmt.Sprintf("%T", o.d.Transactions)).WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions)))
		log.Info("persisting transactions")
		err = o.d.Transactions.ImportTransactions(o.ctx, txns)
		if err != nil {
			log.WithError(err).Error("error saving transactions")
			return 1
		}
		log.Info("persisted transactions")
	}
	return 0
}

func (o ofxImportCommand) Synopsis() string {
	return "Load in downloaded transactions in OFX format."
}
