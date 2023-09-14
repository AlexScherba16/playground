package runner

import (
	log "github.com/sirupsen/logrus"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sort"
	"sync"
)

// postProcessorRunner represents a postprocessor backed by a postprocessing strategy
type postProcessorRunner struct {
	wg               *sync.WaitGroup
	predictorCh      t.PredictorChannel
	postProcessorCh  t.PostProcessorChannel
	postProcStrategy t.PostProcessorStrategy
}

// NewPostProcessorRunner initializes and returns postProcessorRunner
// Returns error if some of wg, predictorCh, postProcessorCh, postProcStrategy is nil
func NewPostProcessorRunner(
	wg *sync.WaitGroup,
	predictorCh t.PredictorChannel,
	postProcessorCh t.PostProcessorChannel,
	postProcStrategy t.PostProcessorStrategy) (*postProcessorRunner, error) {

	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if predictorCh == nil {
		return nil, cerror.NewCustomError("invalid predictor channel")
	}
	if postProcessorCh == nil {
		return nil, cerror.NewCustomError("invalid postprocessor channel")
	}
	if postProcStrategy == nil {
		return nil, cerror.NewCustomError("invalid postprocessor strategy")
	}

	return &postProcessorRunner{
		wg:               wg,
		predictorCh:      predictorCh,
		postProcessorCh:  postProcessorCh,
		postProcStrategy: postProcStrategy,
	}, nil
}

// Run interface implementation, related postprocessing strategy
func (r *postProcessorRunner) Run() {
	defer close(r.postProcessorCh)
	defer r.wg.Done()

	predictions := make([]*t.PredictedData, 0)

	// Read and store predicted data
	for predictData := range r.predictorCh {
		// Received cancel event
		if predictData == nil {
			log.Warning("postprocessor runner shutdown")
			return
		}
		predictions = append(predictions, predictData)
	}

	// Sort predicted data in decreasing order
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Predicted() > predictions[j].Predicted()
	})

	// Convert predicted data to output string, according to postprocessor strategy implementation
	for _, prediction := range predictions {
		r.postProcessorCh <- r.postProcStrategy(prediction)
	}
	log.Debug("postprocessor runner finished work")
}
