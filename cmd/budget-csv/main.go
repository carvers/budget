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
	"github.com/carvers/budget/csv"
	"github.com/carvers/budget/storers"
)

func main() {
	// Set up our logger
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	d := budget.Dependencies{
		Log: yall.New(colour.New(os.Stdout, yall.Severity(level))),
	}

	// Open the file they specified
	if len(os.Args) < 3 {
		d.Log.Error("Usage: budget-csv ACCOUNT_ID /PATH/TO/FILE")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[2])
	if err != nil {
		d.Log.WithError(err).WithField("file", os.Args[2]).Error("error opening file")
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

	pg := storers.NewPostgres(db)
	d.Transactions = pg
	d.Accounts = pg

	ctx := yall.InContext(context.Background(), d.Log)

	transactions, err := csv.FromReader(file, os.Args[1])
	if err != nil {
		d.Log.WithError(err).Error("error parsing Simple transactions")
		os.Exit(1)
	}
	if len(transactions) < 1 {
		d.Log.Warn("no transactions")
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
		d.Log.WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisting transactions")

		// upsert those transactions into postgres
		err = d.Transactions.ImportTransactions(ctx, txns)
		if err != nil {
			d.Log.WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
				WithError(err).Error("error saving transactions")
			os.Exit(1)
		}
		d.Log.WithField("account", os.Args[1]).WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("num_transactions", fmt.Sprintf("%d-%d/%d", i*50, i*50+len(txns), len(transactions))).
			Info("persisted transactions")
	}
}
