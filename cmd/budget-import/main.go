package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	yall "yall.in"
	"yall.in/colour"

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
	log := yall.New(colour.New(os.Stdout, yall.Severity(level)))
	d := budget.Dependencies{
		Log: log,
	}
	ctx := yall.InContext(context.Background(), log)

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		log.WithError(errors.New("no connection string")).
			Error("Error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.WithError(err).Error("Error connecting to Postgres")
		os.Exit(1)
	}

	pg := storers.NewPostgres(db)
	d.Transactions = pg
	d.Accounts = pg
	d.Recurring = pg

	// Set up vault connection
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		log.WithError(errors.New("no vault address")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		log.WithError(errors.New("no vault token")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	d.AccountsSensitive, err = storers.NewVault(vaultAddr, "budget/secret", vaultToken)
	if err != nil {
		log.WithError(err).Error("Error setting up Vault")
		os.Exit(1)
	}

	// get our accounts
	log.WithField("storer", fmt.Sprintf("%T", d.Accounts)).Info("fetching accounts")
	accounts, err := d.Accounts.ListAccounts(ctx)
	if err != nil {
		log.WithError(err).Error("error listing accounts")
		os.Exit(1)
	}

	for _, account := range accounts {
		if !account.Sync {
			log.WithField("account", account.ID).WithField("name", account.Name).
				Info("not configured to sync, ignoring")
			continue
		}
		// retrieve the sensitive info for the account from Vault
		log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.AccountsSensitive)).
			Info("fetching account's sensitive details")
		asd, err := d.AccountsSensitive.GetAccountSensitiveDetails(ctx, account.ID)
		if err != nil {
			log.WithField("account", account.ID).WithError(err).Error("error retrieving sensitive details")
			continue
		}

		// use OFX to retrieve the transactions for this account
		log.WithField("account", account.ID).Info("fetching account transactions")
		transactions, err := ofx.FetchTransactions(ctx, d, account, asd, "3b56383c-2cb6-4c01-8ce0-4951caaf4fa5")
		if err != nil {
			log.WithField("account", account.ID).WithError(err).Error("error fetching transactions")
			continue
		}

		// set our account ID on the transactions
		for pos, txn := range transactions {
			txn.AccountID = account.ID
			transactions[pos] = txn
		}
		log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", len(transactions)).
			Info("Persisting transactions")

		// upsert those transactions into postgres
		err = d.Transactions.ImportTransactions(ctx, transactions)
		if err != nil {
			log.WithError(err).Error("Error saving transactions")
			os.Exit(1)
		}
		log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", len(transactions)).
			Info("Persisted transactions")
	}
}
