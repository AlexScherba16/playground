package aggregator

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/runners/aggregator/campaign"
	"playground/internal/runners/aggregator/country"
	"playground/internal/runners/common"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// NewRunner creates a new data aggregator runner to aggregate records
// According to aggregator parameter
func NewRunner(
	wg *sync.WaitGroup,
	aggregate string,
	recordCh t.RecordChannel,
	aggregateCh t.AggregatorChannel) (common.IRunner, error) {

	// General Factory logic, create data aggregator according to aggregate parameter
	switch aggregate {
	case cnst.AggregateCampaign:
		return campaign.NewAggregator(wg, aggregate, recordCh, aggregateCh)
	case cnst.AggregateCountry:
		return country.NewAggregator(wg, aggregate, recordCh, aggregateCh)
	default:
		return nil, cerror.NewCustomError(fmt.Sprintf("%q invalid aggregate parameter", aggregate))
	}
}
