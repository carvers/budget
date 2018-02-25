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
	"github.com/carvers/budget/simple"
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

	transactions, err := simple.FromReader(ctx, file)
	if err != nil {
		log.WithError(err).Error("error parsing Simple transactions")
		os.Exit(1)
	}
	if len(transactions) < 1 {
		log.Warn("no transactions")
		os.Exit(1)
	}
	accountID := transactions[0].AccountID
	log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Accounts)).
		Info("persisting account")

	err = d.Accounts.CreateAccount(ctx, budget.Account{
		ID:          accountID,
		AccountType: "CHECKING",
		Sync:        false,
	})
	log.WithField("account", accountID).WithField("storer", fmt.Sprintf("%T", d.Accounts)).
		Info("persisted account")

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
