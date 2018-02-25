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

	// Open the file they specified
	if len(os.Args) < 2 {
		log.Error("no filename specified")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.WithError(err).WithField("file", os.Args[1]).Error("error opening file")
		os.Exit(1)
	}

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		log.WithError(errors.New("no connection string")).
			Error("error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.WithError(err).Error("error connecting to Postgres")
		os.Exit(1)
	}

	pg := storers.NewPostgres(db)
	d.Transactions = pg
	d.Accounts = pg

	// Set up vault connection
	vaultAddr := os.Getenv("VAULT_ADDR")
	if postgres == "" {
		log.WithError(errors.New("no vault address")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	vaultToken := os.Getenv("VAULT_TOKEN")
	if postgres == "" {
		log.WithError(errors.New("no vault token")).
			Error("Error setting up Vault")
		os.Exit(1)
	}
	d.AccountsSensitive, err = storers.NewVault(vaultAddr, "budget/secret", vaultToken)
	if err != nil {
		log.WithError(err).Error("Error setting up Vault")
		os.Exit(1)
	}

	// parse our file
	asd, transactions, err := ofx.FromReader(ctx, file)
	if err != nil {
		log.WithError(err).Error("error parsing OFX transactions")
		os.Exit(1)
	}
	if len(transactions) < 1 {
		log.Warn("no transactions")
		os.Exit(1)
	}

	// get a list of accounts so we can list their sensitive info
	log.WithField("storer", fmt.Sprintf("%T", d.Accounts)).Info("fetching accounts")
	accounts, err := d.Accounts.ListAccounts(ctx)
	if err != nil {
		log.WithError(err).Error("error listing accounts")
		os.Exit(1)
	}

	var accountID string

	// get our sensitive details to find out which account this belongs to
	for _, account := range accounts {
		// retrieve the sensitive info for the account from Vault
		log.WithField("account", account.ID).WithField("storer", fmt.Sprintf("%T", d.AccountsSensitive)).
			Info("fetching account's sensitive details")
		a, err := d.AccountsSensitive.GetAccountSensitiveDetails(ctx, account.ID)
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
		log.Error("unable to find account for this file")
		os.Exit(1)
	}

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
		log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisting transactions")

		// upsert those transactions into postgres
		err = d.Transactions.ImportTransactions(ctx, txns)
		if err != nil {
			log.WithError(err).Error("error saving transactions")
			os.Exit(1)
		}
		log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisted transactions")
	}
}
