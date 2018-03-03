package main

import (
	"context"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
)

func serverCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return serverCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type serverCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (s serverCommand) Help() string {
	return "help text"
}

func (s serverCommand) Run(args []string) int {
	return 0
}

func (s serverCommand) Synopsis() string {
	return "Start an API server for access to budget data."
}
