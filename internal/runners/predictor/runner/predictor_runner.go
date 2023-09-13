package runner

import (
	log "github.com/sirupsen/logrus"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// inputWorkerChanMap map to store input aggregated data channels for
// Each unique predictor worker
type inputWorkerChanMap map[string]t.AggregatorChannel

// predictorRunner represents a data predictor backed by a prediction strategy
type predictorRunner struct {
	wg           *sync.WaitGroup
	aggregatorCh t.AggregatorChannel
	predictorCh  t.PredictorChannel
	prStrategy   t.PredictWorkerStrategy
}

// NewPredictorRunner initializes and returns predictorRunner
// Returns error if some of wg, aggregatorCh, predictorCh, prStrategy is nil
func NewPredictorRunner(
	wg *sync.WaitGroup,
	aggregatorCh t.AggregatorChannel,
	predictorCh t.PredictorChannel,
	prStrategy t.PredictWorkerStrategy) (*predictorRunner, error) {

	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if aggregatorCh == nil {
		return nil, cerror.NewCustomError("invalid aggregator channel")
	}
	if predictorCh == nil {
		return nil, cerror.NewCustomError("invalid predictor channel")
	}
	if prStrategy == nil {
		return nil, cerror.NewCustomError("invalid predictor strategy worker")
	}

	return &predictorRunner{
		wg:           wg,
		aggregatorCh: aggregatorCh,
		predictorCh:  predictorCh,
		prStrategy:   prStrategy,
	}, nil
}

// Run interface implementation, related strategy prediction model
func (r *predictorRunner) Run() {
	defer close(r.predictorCh)
	defer r.wg.Done()

	// Prepare worker sync and communication attributes
	workerOutCh := t.NewPredictorChannel(5)
	workerInChannelMap := make(inputWorkerChanMap)
	workerWg := &sync.WaitGroup{}

	// Read aggregated data until aggregate channel is open
	for aggData := range r.aggregatorCh {
		// Received cancel event
		if aggData == nil {
			log.Warning("predict runner shutdown")
			r.predictorCh <- nil

			// Shutdown all running workers
			for _, channel := range workerInChannelMap {
				channel <- nil
				close(channel)
			}

			// Wait until workers stop running
			workerWg.Wait()
			log.Warning("predict workers released")
			return
		}

		// Spinup new worker in case of unique aggregated data received
		if _, found := workerInChannelMap[aggData.Key()]; !found {
			workerInCh := t.NewAggregatorChannel(2)
			workerWg.Add(1)
			go r.prStrategy(workerWg, aggData.Key(), workerInCh, workerOutCh)
			workerInChannelMap[aggData.Key()] = workerInCh
		}

		// Send aggregated data to key related worker
		workerInChannelMap[aggData.Key()] <- aggData
	}

	// Aggregated data channel closed
	// Get workers result and close all workers channels
	for _, workerInputChannel := range workerInChannelMap {
		// Release goroutines
		close(workerInputChannel)
		tmp := <-workerOutCh
		r.predictorCh <- tmp
	}

	// Wait until workers stop running
	workerWg.Wait()
	log.Debug("predict runner finished work")
}
