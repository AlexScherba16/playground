package types

import "sync"

// AggregatorStrategy strategy for Record data aggregation algorithm
type AggregatorStrategy func(record *Record) *AggregatedData

// PredictWorkerStrategy strategy for data prediction algorithm
type PredictWorkerStrategy func(
	wg *sync.WaitGroup,
	key string,
	inCh AggregatorChannel,
	outCh PredictorChannel)

// PostProcessorStrategy strategy for PredictedData to string conversion algorithm
type PostProcessorStrategy func(predictedData *PredictedData) string
