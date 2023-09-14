package aggregator_factory

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/runners/aggregator/runner"
	campaign "playground/internal/runners/aggregator/strategy/campaign"
	country "playground/internal/runners/aggregator/strategy/country"
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
		return runner.NewAggregatorRunner(wg, recordCh, aggregateCh, campaign.NewCampaignAggregatorStrategy())
	case cnst.AggregateCountry:
		return runner.NewAggregatorRunner(wg, recordCh, aggregateCh, country.NewCountryAggregatorStrategy())
	default:
		return nil, cerror.NewCustomError(fmt.Sprintf("%q invalid aggregate parameter", aggregate))
	}
}
