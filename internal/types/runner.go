package types

import "sync"

type PredictWorkerStrategy func(
	wg *sync.WaitGroup,
	key string,
	inCh AggregatorChannel,
	outCh PredictorChannel)

// PostProcessorStrategy strategy for PredictedData to string conversion algorithm
type PostProcessorStrategy func(predictedData *PredictedData) string
