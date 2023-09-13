package types

import "sync"

type PredictWorkerStrategy func(
	wg *sync.WaitGroup,
	key string,
	inCh AggregatorChannel,
	outCh PredictorChannel)
