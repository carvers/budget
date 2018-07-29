package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"darlinggo.co/version"
	"github.com/carvers/budget"
	"github.com/carvers/budget/storers"
	"github.com/mitchellh/cli"
	yall "yall.in"
	"yall.in/colour"
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
	var db *sql.DB
	var err error
	if postgres != "" {
		db, err = sql.Open("postgres", postgres)
		if err != nil {
			log.WithError(err).Error("Error connecting to Postgres")
			os.Exit(1)
		}

		pg := storers.NewPostgres(db)
		d.Transactions = pg
		d.Accounts = pg
		d.Groups = pg
	}

	// Set up vault connection
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken != "" && vaultAddr != "" {
		d.AccountsSensitive, err = storers.NewVault(vaultAddr, "budget/secret", vaultToken)
		if err != nil {
			log.WithError(err).Error("Error setting up Vault")
			os.Exit(1)
		}
	}

	c := cli.NewCLI("budgetctl", fmt.Sprintf("%s (%s)", version.Tag, version.Hash))
	c.Args = os.Args[1:]
	ui := &cli.ColoredUi{
		InfoColor:  cli.UiColorCyan,
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
	c.Commands = map[string]cli.CommandFactory{
		// fetch ofx transactions
		"ofx fetch": ofxFetchCommandFactory(ctx, ui, d),

		// import ofx transactions
		"ofx import": ofxImportCommandFactory(ctx, ui, d),

		// detect groups
		"groups detect": groupsDetectCommandFactory(ctx, ui, d),
		// detect periods
		"groups periods": groupsPeriodsCommandFactory(ctx, ui, d),

		// import csv transactions
		"csv import": csvImportCommandFactory(ctx, ui, d),

		// import simple transactions
		"simple import": simpleImportCommandFactory(ctx, ui, d),

		// run the API server
		"server": serverCommandFactory(ctx, ui, d),

		// run database migrations
		"migrations run": migrationsRunCommandFactory(ctx, ui, db),
	}
	c.Autocomplete = true

	status, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}
