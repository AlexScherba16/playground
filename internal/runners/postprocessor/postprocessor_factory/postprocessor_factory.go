package postprocessor_factory

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/runners/common"
	"playground/internal/runners/postprocessor/runner"
	"playground/internal/runners/postprocessor/strategy/campaign_postprocessor_strategy"
	"playground/internal/runners/postprocessor/strategy/country_postprocessor_strategy"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// NewRunner creates a new data postprocessor runner to prepare predicted data for output
// According to aggregate parameter
func NewRunner(
	wg *sync.WaitGroup,
	aggregate string,
	predictCh t.PredictorChannel,
	postCh t.PostProcessorChannel) (common.IRunner, error) {

	// General Factory logic, create data predictor according to aggregate parameter
	switch aggregate {
	case cnst.AggregateCountry:
		return runner.NewPostProcessorRunner(wg, predictCh, postCh,
			country_postprocessor_strategy.NewPostProcessorStrategy())
	case cnst.AggregateCampaign:
		return runner.NewPostProcessorRunner(wg, predictCh, postCh,
			campaign_postprocessor_strategy.NewPostProcessorStrategy())
	default:
		return nil, cerror.NewCustomError(fmt.Sprintf("%q invalid postprocessor parameter", aggregate))
	}
}
