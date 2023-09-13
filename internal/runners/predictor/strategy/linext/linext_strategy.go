package linext

import (
	log "github.com/sirupsen/logrus"
	cnst "playground/internal/constants"
	t "playground/internal/types"
	"playground/internal/utils/predictor"
	"sync"
)

// linearExtrapolationWorker perform linear extrapolation prediction logic for key related aggregated data
// IMPORTANT: it doesn't close channels and predict only for PredictForNDay
func linearExtrapolationWorker(wg *sync.WaitGroup, key string, inCh t.AggregatorChannel, outCh t.PredictorChannel) {
	defer wg.Done()

	ltvSums := t.LtvCollection{}
	ltvNonEmptyValues := make([]int, cnst.LtvLen)

	// Read aggregated data
	for aggData := range inCh {
		// Received cancel event
		if aggData == nil {
			log.Warning("linext worker shutdown")
			return
		}

		// Collect ltvData, calculate non zero values
		for i, value := range aggData.Ltv() {
			// Scip 0 values
			if value == 0 {
				continue
			}
			ltvSums[i] += value
			ltvNonEmptyValues[i]++
		}
	}

	// All ltvData collected here, calculate average and store non zero values
	averages := make([]float64, 0)
	for i, sum := range ltvSums {
		if ltvNonEmptyValues[i] != 0 {
			averages = append(averages, sum/float64(ltvNonEmptyValues[i]))
		}
	}
	// Predict n-th day ltv
	result := t.NewPredictedData(key, predictor.LinearExtrapolation(averages, cnst.PredictForNDay))

	// Send n-th day predicted data
	outCh <- result
}

func NewPredictWorkerStrategy() t.PredictWorkerStrategy {
	return linearExtrapolationWorker
}
