package nbastats

import "time"

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// PerMode determines if the stats to be queried should be the per game average or the cumulative totals
	PerMode string `json:"per_mode" validate:"required" default:"PerGame"`
	// how often the connector will get data from the url
	PollingPeriod time.Duration `json:"pollingPeriod" default:"5m"`
}
