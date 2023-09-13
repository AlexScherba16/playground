package country

import (
	"fmt"
	cnst "playground/internal/constants"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"

	log "github.com/sirupsen/logrus"
)

// countryAggregator represents a data aggregator backed by a "country" aggregate parameter
type countryAggregator struct {
	wg           *sync.WaitGroup
	aggregate    string
	recordCh     t.RecordChannel
	aggregatedCh t.AggregateChannel
}

// NewAggregator initializes and returns countryAggregator
// Returns error if some of wg, recordCh, aggregatedCh is nil, or invalid aggregate parameter
func NewAggregator(
	wg *sync.WaitGroup,
	aggregate string,
	recordCh t.RecordChannel,
	aggregateCh t.AggregateChannel) (*countryAggregator, error) {

	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if aggregate != cnst.AggregateCountry {
		return nil, cerror.NewCustomError(fmt.Sprintf("%q should be provided as an aggregate parameter, not %q",
			cnst.AggregateCountry, aggregate))
	}
	if recordCh == nil {
		return nil, cerror.NewCustomError("invalid record channel")
	}
	if aggregateCh == nil {
		return nil, cerror.NewCustomError("invalid aggregate channel")
	}

	return &countryAggregator{
		wg:           wg,
		aggregate:    aggregate,
		recordCh:     recordCh,
		aggregatedCh: aggregateCh,
	}, nil
}

// Run interface implementation, related to country aggregation logic
func (r *countryAggregator) Run() {
	defer close(r.aggregatedCh)
	defer r.wg.Done()

	// Read records until record channel is open
	for record := range r.recordCh {
		// Received cancel event
		if record == nil {
			log.Warning("campaign aggregator shutdown")

			// Notify next runner about cancel event
			r.aggregatedCh <- nil
			return
		}

		// Send aggregated data to next runner
		aggregatedData := t.NewAggregatedData(record.Country(), record.Ltv())
		r.aggregatedCh <- aggregatedData
	}
	log.Debug("campaign aggregator finished work")
}
