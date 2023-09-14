package country

import (
	t "playground/internal/types"
)

// countryAggregatorStrategy Record aggregation strategy, accorded to country key
func countryAggregatorStrategy(record *t.Record) *t.AggregatedData {
	return t.NewAggregatedData(record.Country(), record.Ltv())
}

// NewCountryAggregatorStrategy returns country aggregator strategy
func NewCountryAggregatorStrategy() t.AggregatorStrategy {
	return countryAggregatorStrategy
}
