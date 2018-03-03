package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/carvers/budget/migrations"
	"github.com/mitchellh/cli"
	migrate "github.com/rubenv/sql-migrate"
	yall "yall.in"
)

func migrationsRunCommandFactory(ctx context.Context, ui cli.Ui, db *sql.DB) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return migrationsRunCommand{
			db:  db,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type migrationsRunCommand struct {
	db  *sql.DB
	ui  cli.Ui
	ctx context.Context
}

func (m migrationsRunCommand) Help() string {
	return "help text"
}

func (m migrationsRunCommand) Run(args []string) int {
	log := yall.FromContext(m.ctx).WithField("db", fmt.Sprintf("%T", m.db))
	if m.db == nil {
		log.Error("must have database configured to run migrations on")
		return 1
	}
	migrations := &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "sql",
	}
	log.Info("running migrations")
	_, err := migrate.Exec(m.db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithError(err).Error("error running migrations")
		return 1
	}
	log.Info("migrations completed")
	return 0
}

func (m migrationsRunCommand) Synopsis() string {
	return "Run database schema migrations."
}
