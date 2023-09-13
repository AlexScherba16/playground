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

// AggregateChannel is a channel type for transmitting AggregatedData instances
type AggregateChannel chan *AggregatedData

// NewAggregateChannel initializes and returns AggregateChannel with a buffer size
func NewAggregateChannel(aggregatedBuffer uint) AggregateChannel {
	return make(AggregateChannel, aggregatedBuffer)
}

// Channels is channel container util
type Channels struct {
	RecordCh    RecordChannel
	ErrorCh     ErrorChannel
	AggregateCh AggregateChannel
}

// NewChannels initializes and returns a Channels structure
func NewChannels(recordsBuffer, errorsBuffer, aggregateBuffer uint) *Channels {
	return &Channels{
		RecordCh:    NewRecordChannel(recordsBuffer),
		ErrorCh:     NewErrorChannel(errorsBuffer),
		AggregateCh: NewAggregateChannel(aggregateBuffer),
	}
}
