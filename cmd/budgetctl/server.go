package main

import (
	"context"
	"net/http"

	"github.com/carvers/budget"
	"github.com/carvers/budget/apiv1"
	"github.com/mitchellh/cli"
	yall "yall.in"
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
	v1 := apiv1.APIv1{
		s.d,
	}
	http.Handle("/api/v1/", v1.Server("/api/v1/"))
	yall.FromContext(s.ctx).WithField("address", ":9323").Info("API server started")
	err := http.ListenAndServe(":9323", nil)
	if err != nil {
		yall.FromContext(s.ctx).WithError(err).Error("error starting server")
		return 1
	}
	return 0
}

func (s serverCommand) Synopsis() string {
	return "Start an API server for access to budget data."
}
