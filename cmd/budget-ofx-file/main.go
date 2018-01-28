package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/pkg/errors"

	"github.com/carvers/budget"
	"github.com/carvers/budget/ofx"
	"github.com/carvers/budget/storers"
)

func main() {
	// Set up our logger
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	lvl, err := log.ParseLevel(level)
	if err != nil {
		fmt.Printf("Invalid log level %q; please set LOG_LEVEL to %q, %q, %q, %q, or %q.",
			log.DebugLevel.String(), log.InfoLevel.String(), log.WarnLevel.String(),
			log.ErrorLevel.String(), log.FatalLevel.String())
		os.Exit(1)
	}
	d := budget.Dependencies{
		Log: &log.Logger{
			Level:   lvl,
			Handler: cli.Default,
		},
	}

	// Open the file they specified
	if len(os.Args) < 2 {
		d.Log.Error("no filename specified")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		d.Log.WithError(err).WithField("file", os.Args[1]).Error("error opening file")
		os.Exit(1)
	}

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		d.Log.WithError(errors.New("no connection string")).
			Error("error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		d.Log.WithError(err).Error("error connecting to Postgres")
		os.Exit(1)
	}

	pg := storers.NewPostgres(db, d.Log)
	d.Transactions = pg
	d.Accounts = pg

	// Set up vault connection
	vaultAddr := os.Getenv("VAULT_ADDR")
	if postgres == "" {
		d.Log.WithError(errors.New("no vault address")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	vaultToken := os.Getenv("VAULT_TOKEN")
	if postgres == "" {
		d.Log.WithError(errors.New("no vault token")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	d.AccountsSensitive, err = storers.NewVault(vaultAddr, "budget/secret", d.Log, vaultToken)
	if err != nil {
		d.Log.WithError(err).Error("Error setting up Vault")
		os.Exit(1)
	}

	// parse our file
	asd, transactions, err := ofx.FromReader(file)
	if err != nil {
		d.Log.WithError(err).Error("error parsing OFX transactions")
		os.Exit(1)
	}
	if len(transactions) < 1 {
		d.Log.Warn("no transactions")
		os.Exit(1)
	}

	// get a list of accounts so we can list their sensitive info
	d.Log.WithField("storer", fmt.Sprintf("%T", d.Accounts)).Info("fetching accounts")
	accounts, err := d.Accounts.ListAccounts()
	if err != nil {
		d.Log.WithError(err).Error("error listing accounts")
		os.Exit(1)
	}

	var accountID string

	// get our sensitive details to find out which account this belongs to
	for _, account := range accounts {
		// retrieve the sensitive info for the account from Vault
		d.Log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.AccountsSensitive)).
			Info("fetching account's sensitive details")
		a, err := d.AccountsSensitive.GetAccountSensitiveDetails(account.ID)
		if err != nil {
			d.Log.WithField("account", account.ID).WithError(err).Error("error retrieving sensitive details")
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
		d.Log.Error("unable to find account for this file")
		os.Exit(1)
	}

	for i := 0; i < (len(transactions)/50)+1; i++ {
		if len(transactions) <= i*50 {
			continue
		}
		txns := transactions[i*50:]
		if len(txns) > 50 {
			txns = txns[:50]
		}
		d.Log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisting transactions")

		// upsert those transactions into postgres
		err = d.Transactions.ImportTransactions(txns)
		if err != nil {
			d.Log.WithError(err).Error("error saving transactions")
			os.Exit(1)
		}
		d.Log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisted transactions")
	}
}
