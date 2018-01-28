package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/carvers/budget"
	"github.com/carvers/budget/migrations"
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
		d.Log.WithError(errors.New("no connection string set")).Error("Error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		d.Log.WithError(err).Error("Error connecting to Postgres")
		os.Exit(1)
	}

	// run our postgres migrations
	migrations := &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "sql",
	}
	d.Log.Info("Running migrations")
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		d.Log.WithError(err).Error("Error running migrations for Postgres")
		os.Exit(1)
	}
	d.Log.Info("Migrations completed.")
}
