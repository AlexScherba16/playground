package types

// RecordChannel is a channel type for transmitting Record instances
type RecordChannel chan *Record

// NewRecordChannel initializes and returns RecordChannel with a buffer size
func NewRecordChannel(recordsBuffer uint) RecordChannel {
	return make(RecordChannel, recordsBuffer)
}

// ErrorChannel is a channel type for transmitting error instances
type ErrorChannel chan error

// NewErrorChannel initializes and returns ErrorChannel with a buffer size
func NewErrorChannel(errorsBuffer uint) ErrorChannel {
	return make(ErrorChannel, errorsBuffer)
}

// AggregatorChannel is a channel type for transmitting AggregatedData instances
type AggregatorChannel chan *AggregatedData

// NewAggregatorChannel initializes and returns AggregatorChannel with a buffer size
func NewAggregatorChannel(aggregatedBuffer uint) AggregatorChannel {
	return make(AggregatorChannel, aggregatedBuffer)
}

// PredictorChannel is a channel type for transmitting PredictedData instances
type PredictorChannel chan *PredictedData

// NewPredictorChannel initializes and returns PredictorChannel with a buffer size
func NewPredictorChannel(predictBuffer uint) PredictorChannel {
	return make(PredictorChannel, predictBuffer)
}

// Channels is channel container util
type Channels struct {
	RecordCh    RecordChannel
	ErrorCh     ErrorChannel
	AggregateCh AggregatorChannel
	PredictCh   PredictorChannel
}

// NewChannels initializes and returns a Channels structure
func NewChannels(recordsBuffer, errorsBuffer, aggregateBuffer, predictBuffer uint) *Channels {
	return &Channels{
		RecordCh:    NewRecordChannel(recordsBuffer),
		ErrorCh:     NewErrorChannel(errorsBuffer),
		AggregateCh: NewAggregatorChannel(aggregateBuffer),
		PredictCh:   NewPredictorChannel(predictBuffer),
	}
}
