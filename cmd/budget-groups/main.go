package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/pkg/errors"

	"github.com/carvers/budget"
	"github.com/carvers/budget/storers"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	lvl, err := log.ParseLevel(level)
	if err != nil {
		fmt.Printf("Invalid log level %q; please set LOG_LEVEL to %q, %q, %q, %q, or %q.",
			level, log.DebugLevel.String(), log.InfoLevel.String(), log.WarnLevel.String(),
			log.ErrorLevel.String(), log.FatalLevel.String())
		os.Exit(1)
	}
	d := budget.Dependencies{
		// Set up our logger
		Log: &log.Logger{
			Level:   lvl,
			Handler: cli.Default,
		},
	}

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		d.Log.WithError(errors.New("no connection string set")).Error("error setting up Postgres")
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
	d.Recurring = pg

	groups, err := budget.GroupTransactions(d)
	if err != nil {
		d.Log.WithError(err).Error("error finding groups")
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

	d.Log.WithField("num_groups", len(recurs)).WithField("storer", fmt.Sprintf("%T", d.Recurring)).
		Info("storing recurring groups")
	err = d.Recurring.CreateRecurrings(recurs)
	if err != nil {
		d.Log.WithField("storer", fmt.Sprintf("%T", d.Recurring)).WithError(err).
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
		d.Log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).
			WithField("recurring_id", recurID).
			Info("updating RecurringID on transactions")
		err = d.Transactions.UpdateTransactions(tf, change)
		if err != nil {
			d.Log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).
				WithField("recurring_id", recurID).
				Error("error updating RecurringID on transactions")
		}
	}
}
