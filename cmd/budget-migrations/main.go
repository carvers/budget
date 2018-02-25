package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	yall "yall.in"
	"yall.in/colour"

	"github.com/carvers/budget/migrations"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	log := yall.New(colour.New(os.Stdout, yall.Severity(level)))

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		log.WithError(errors.New("no connection string set")).Error("Error setting up Postgres")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.WithError(err).Error("Error connecting to Postgres")
		os.Exit(1)
	}

	// run our postgres migrations
	migrations := &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "sql",
	}
	log.Info("Running migrations")
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithError(err).Error("Error running migrations for Postgres")
		os.Exit(1)
	}
	log.Info("Migrations completed.")
}
