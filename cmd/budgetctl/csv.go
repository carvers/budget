package main

import (
	"context"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
)

func csvImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return csvImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type csvImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (c csvImportCommand) Help() string {
	return "help text"
}

func (c csvImportCommand) Run(args []string) int {
	return 0
}

func (c csvImportCommand) Synopsis() string {
	return "Load in downloaded transactions in CSV format."
}
