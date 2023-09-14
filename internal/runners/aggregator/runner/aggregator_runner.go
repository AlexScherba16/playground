package runner

import (
	log "github.com/sirupsen/logrus"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// countryAggregator represents a data aggregator backed by a "country" aggregate parameter
type countryAggregator struct {
	wg                  *sync.WaitGroup
	recordCh            t.RecordChannel
	aggregatedCh        t.AggregatorChannel
	aggregationStrategy t.AggregatorStrategy
}

// NewAggregatorRunner initializes and returns countryAggregator
// Returns error if some of wg, recordCh, aggregatedCh, aggregationStrategy is nil
func NewAggregatorRunner(
	wg *sync.WaitGroup,
	recordCh t.RecordChannel,
	aggregateCh t.AggregatorChannel,
	aggregationStrategy t.AggregatorStrategy) (*countryAggregator, error) {

	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if recordCh == nil {
		return nil, cerror.NewCustomError("invalid record channel")
	}
	if aggregateCh == nil {
		return nil, cerror.NewCustomError("invalid aggregate channel")
	}
	if aggregationStrategy == nil {
		return nil, cerror.NewCustomError("invalid aggregation strategy")
	}

	return &countryAggregator{
		wg:                  wg,
		recordCh:            recordCh,
		aggregatedCh:        aggregateCh,
		aggregationStrategy: aggregationStrategy,
	}, nil
}

// Run interface implementation, related to data aggregation strategy
func (r *countryAggregator) Run() {
	defer close(r.aggregatedCh)
	defer r.wg.Done()

	// Read records until record channel is open
	for record := range r.recordCh {
		// Received cancel event
		if record == nil {
			log.Warning("aggregator runner shutdown")

			// Notify next runner about cancel event
			r.aggregatedCh <- nil
			return
		}

		// Send aggregated data to next runner
		r.aggregatedCh <- r.aggregationStrategy(record)
	}
	log.Debug("aggregator runner finished work")
}
