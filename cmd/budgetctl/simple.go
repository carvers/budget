package main

import (
	"context"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
)

func simpleImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return simpleImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type simpleImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (s simpleImportCommand) Help() string {
	return "help text"
}

func (s simpleImportCommand) Run(args []string) int {
	return 0
}

func (s simpleImportCommand) Synopsis() string {
	return "Load in downloaded transactions from Simple in JSON format."
}
