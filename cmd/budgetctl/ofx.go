package main

import (
	"context"

	"github.com/carvers/budget"
	"github.com/mitchellh/cli"
)

func ofxFetchCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return ofxFetchCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type ofxFetchCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (o ofxFetchCommand) Help() string {
	return "help text"
}

func (o ofxFetchCommand) Run(args []string) int {
	return 0
}

func (o ofxFetchCommand) Synopsis() string {
	return "Load in recent transactions from your accounts using OFX."
}

func ofxImportCommandFactory(ctx context.Context, ui cli.Ui, d budget.Dependencies) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return ofxImportCommand{
			d:   d,
			ui:  ui,
			ctx: ctx,
		}, nil
	}
}

type ofxImportCommand struct {
	d   budget.Dependencies
	ui  cli.Ui
	ctx context.Context
}

func (o ofxImportCommand) Help() string {
	return "help text"
}

func (o ofxImportCommand) Run(args []string) int {
	return 0
}

func (o ofxImportCommand) Synopsis() string {
	return "Load in downloaded transactions in OFX format."
}
