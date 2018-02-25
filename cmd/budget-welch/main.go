package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	yall "yall.in"
	"yall.in/colour"

	"github.com/mjibson/go-dsp/spectral"
	"github.com/pkg/errors"

	"github.com/carvers/budget"
	"github.com/carvers/budget/storers"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	log := yall.New(colour.New(os.Stdout, yall.Severity(level)))
	d := budget.Dependencies{
		// Set up our logger
		Log: log,
	}
	ctx := yall.InContext(context.Background(), log)

	if len(os.Args) < 2 {
		log.Error("no recurring ID set")
		os.Exit(1)
	}

	// Set up postgres connection
	postgres := os.Getenv("PG_DB")
	if postgres == "" {
		log.WithError(errors.New("no connection string set")).Error("error setting up Postgres")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.WithError(err).Error("error connecting to Postgres")
		os.Exit(1)
	}

	pg := storers.NewPostgres(db)
	d.Transactions = pg

	transactions, err := d.Transactions.ListTransactions(ctx, budget.TransactionFilters{
		RecurringID: &os.Args[1],
	})

	if err != nil {
		log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).WithError(err).
			Error("error retrieving transactions")
		os.Exit(1)
	}
	if len(transactions) < 1 {
		log.WithField("storer", fmt.Sprintf("%T", d.Transactions)).WithField("recurring_id", os.Args[1]).
			Error("no transactions match that recurring_id")
		os.Exit(1)
	}

	// sort our transactions by the day they were posted
	sort.Slice(transactions, func(i, j int) bool { return transactions[i].DatePosted.Before(transactions[j].DatePosted) })

	// a map of transactions by the day they happened on
	days := map[int64][]budget.Transaction{}

	var totalDays int64
	for _, txn := range transactions {
		offset := int64(txn.DatePosted.Sub(transactions[0].DatePosted).Truncate(time.Hour) / time.Hour / 24)
		days[offset] = append(days[offset], txn)
		if offset > totalDays {
			totalDays = offset
		}
	}

	log.WithField("total_days", totalDays).Debug("")

	log.WithField("populated_days", len(days)).Debug("")

	transactionsPerDay := make([]float64, totalDays+1)
	centsPerDay := make([]float64, totalDays+1)
	for day, txns := range days {
		log.WithField("day", day).Debug("")
		transactionsPerDay[day] = float64(len(txns))
		for _, txn := range txns {
			centsPerDay[day] += float64(txn.Amount)
		}
	}

	txnDensities, txnFreqs := spectral.Pwelch(transactionsPerDay, 1.0, &spectral.PwelchOptions{})
	txnPSD := zipDensitiesAndFrequencies(txnDensities, txnFreqs)
	sort.Slice(txnPSD, func(i, j int) bool { return txnPSD[i].density > txnPSD[j].density })
	for _, psd := range txnPSD {
		log.WithField("density", psd.density).WithField("frequency", psd.frequency).Info("calculated PSD for transactions")
	}

	centDensities, centFreqs := spectral.Pwelch(centsPerDay, 1.0, &spectral.PwelchOptions{})
	centPSD := zipDensitiesAndFrequencies(centDensities, centFreqs)
	sort.Slice(centPSD, func(i, j int) bool { return centPSD[i].density > centPSD[j].density })
	for _, psd := range centPSD {
		log.WithField("density", psd.density).WithField("frequency", psd.frequency).Info("calculated PSD for amount")
	}

	p, err := plot.New()
	if err != nil {
		log.WithError(err).Error("error generating plot")
		os.Exit(1)
	}
	p.Title.Text = "Transactions per Day"
	p.X.Label.Text = "Day"
	p.Y.Label.Text = "# of Transactions"

	xy := make(plotter.XYs, len(transactionsPerDay))
	for day, count := range transactionsPerDay {
		log.WithField("day", day).WithField("count", count).Debug("adding point")
		xy[day].X = float64(day)
		xy[day].Y = float64(count)
	}

	err = plotutil.AddLinePoints(p, xy)
	if err != nil {
		log.WithError(err).Error("error adding points")
		os.Exit(1)
	}

	err = p.Save(12*vg.Inch, 12*vg.Inch, "/usr/local/var/www/txnsperday.png")
	if err != nil {
		log.WithError(err).Error("error writing plot to png")
		os.Exit(1)
	}
}

type psd struct {
	density   float64
	frequency float64
}

func zipDensitiesAndFrequencies(d, f []float64) []psd {
	res := make([]psd, 0, len(d))
	for pos, den := range d {
		res = append(res, psd{
			density:   den,
			frequency: f[pos],
		})
	}
	return res
}
