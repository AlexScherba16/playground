package predictor_factory

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/runners/common"
	pr "playground/internal/runners/predictor/runner"
	"playground/internal/runners/predictor/strategy/linext"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// NewRunner creates a new data predictor runner to perform predictions on aggregated data
// According to model parameter
func NewRunner(
	wg *sync.WaitGroup,
	model string,
	aggregateCh t.AggregatorChannel,
	predictCh t.PredictorChannel) (common.IRunner, error) {

	// General Factory logic, create data predictor according to model parameter
	switch model {
	case cnst.LinearExtrapolationPredictorModel:
		return pr.NewPredictorRunner(wg, aggregateCh, predictCh, linext.NewPredictWorkerStrategy())
	default:
		return nil, cerror.NewCustomError(fmt.Sprintf("%q invalid model parameter", model))
	}
}
