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
	"github.com/carvers/budget/storers"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	log := yall.New(colour.New(os.Stdout, yall.Severity(level)))
	d := budget.Dependencies{
		// Set up our logger
		Log: log,
	}

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		log.WithError(errors.New("no connection string set")).Error("error setting up Postgres")
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
	d.Recurring = pg

	ctx := yall.InContext(context.Background(), log)

	groups, err := budget.GroupTransactions(ctx, d)
	if err != nil {
		log.WithError(err).Error("error finding groups")
		os.Exit(1)
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

	log.WithField("num_groups", len(recurs)).WithField("storer", fmt.Sprintf("%T", d.Recurring)).
		Info("storing recurring groups")
	err = d.Recurring.CreateRecurrings(ctx, recurs)
	if err != nil {
		log.WithField("storer", fmt.Sprintf("%T", d.Recurring)).WithError(err).
			Error("error storing recurring groups")
		os.Exit(1)
	}

	for recur, txns := range txnChanges {
		recurID := recur
		tf := budget.TransactionFilters{
			IDs: txns,
		}
		change := budget.TransactionChange{
			RecurringID: &recurID,
		}
		log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("recurring_id", recurID).
			Info("updating RecurringID on transactions")
		err = d.Transactions.UpdateTransactions(ctx, tf, change)
		if err != nil {
			log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).
				WithField("recurring_id", recurID).
				Error("error updating RecurringID on transactions")
		}
	}
}
