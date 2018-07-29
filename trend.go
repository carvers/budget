package budget

// Trend represents recurring transactions in a category and should work out to a stable pattern. It could be used for
// actual recurring budget items, like a paycheck or monthly subscription, or it could be used to track a pseudo-item,
// like spending on groceries or buying gas or eating out. Trends will be used to forecast future transactions when
// attempting to simulate finances on future dates.
type Trend struct {
	// group transactions across this many days, and lump them in as "one" transaction for statistical analysis
	// for things that may happen across many small transactions but should still work out to a pattern
	// e.g., groceries, gas, etc. For actual, stable, one-transaction-per-period, just leave this at 0, and we'll
	// figure out the period based on the transactions.
	ForcedPeriodDays int

	// a friendly name to identify this trend with
	Name string
}
