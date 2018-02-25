package main

import (
	"context"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
)

func groupsDetectCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return groupsDetectCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type groupsDetectCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (g groupsDetectCommand) Help() string {
	return "help text"
}

func (g groupsDetectCommand) Run(args []string) int {
	return 0
}

func (g groupsDetectCommand) Synopsis() string {
	return "Detect transactions that may belong to a recurring group, and update them to be part of it."
}
