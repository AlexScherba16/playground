package campaign

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	cnst "playground/internal/constants"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// campaignAggregator represents a data aggregator backed by a "campaign" aggregate parameter
type campaignAggregator struct {
	wg           *sync.WaitGroup
	aggregate    string
	recordCh     t.RecordChannel
	aggregatedCh t.AggregateChannel
}

// NewAggregator initializes and returns campaignAggregator
// Returns error if some of wg, recordCh, aggregatedCh is nil, or invalid aggregate parameter
func NewAggregator(
	wg *sync.WaitGroup,
	aggregate string,
	recordCh t.RecordChannel,
	aggregateCh t.AggregateChannel) (*campaignAggregator, error) {

	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if aggregate != cnst.AggregateCampaign {
		return nil, cerror.NewCustomError(fmt.Sprintf("%q should be provided as an aggregate parameter, not %q",
			cnst.AggregateCampaign, aggregate))
	}
	if recordCh == nil {
		return nil, cerror.NewCustomError("invalid record channel")
	}
	if aggregateCh == nil {
		return nil, cerror.NewCustomError("invalid aggregate channel")
	}

	return &campaignAggregator{
		wg:           wg,
		aggregate:    aggregate,
		recordCh:     recordCh,
		aggregatedCh: aggregateCh,
	}, nil
}

// Run interface implementation, related to campaign aggregation logic
func (r *campaignAggregator) Run() {
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
		aggregatedData := t.NewAggregatedData(t.KeyType(record.CampaignId()), record.Ltv())
		r.aggregatedCh <- aggregatedData
	}
	log.Debug("campaign aggregator finished work")
}
