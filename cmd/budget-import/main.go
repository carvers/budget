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

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		d.Log.WithError(errors.New("no connection string")).
			Error("Error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		d.Log.WithError(err).Error("Error connecting to Postgres")
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

	// get our accounts
	d.Log.WithField("storer", fmt.Sprintf("%T", d.Accounts)).Info("fetching accounts")
	accounts, err := d.Accounts.ListAccounts()
	if err != nil {
		d.Log.WithError(err).Error("error listing accounts")
		os.Exit(1)
	}

	for _, account := range accounts {
		// retrieve the sensitive info for the account from Vault
		d.Log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.AccountsSensitive)).
			Info("fetching account's sensitive details")
		asd, err := d.AccountsSensitive.GetAccountSensitiveDetails(account.ID)
		if err != nil {
			d.Log.WithField("account", account.ID).WithError(err).Error("error retrieving sensitive details")
			continue
		}

		// use OFX to retrieve the transactions for this account
		d.Log.WithField("account", account.ID).Info("fetching account transactions")
		transactions, err := ofx.FetchTransactions(d, account, asd, "3b56383c-2cb6-4c01-8ce0-4951caaf4fa5")
		if err != nil {
			d.Log.WithField("account", account.ID).WithError(err).Error("error fetching transactions")
			continue
		}
		d.Log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", len(transactions)).
			Info("Persisting transactions")

		// upsert those transactions into postgres
		err = d.Transactions.ImportTransactions(transactions)
		if err != nil {
			d.Log.WithError(err).Error("Error saving transactions")
			os.Exit(1)
		}
		d.Log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", len(transactions)).
			Info("Persisted transactions")
	}
}
