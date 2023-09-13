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

// Channels is channel container util
type Channels struct {
	RecordCh RecordChannel
	ErrorCh  ErrorChannel
}

// NewChannels initializes and returns a Channels structure
func NewChannels(recordsBuffer, errorsBuffer uint) *Channels {
	return &Channels{
		RecordCh: make(RecordChannel, recordsBuffer),
		ErrorCh:  make(ErrorChannel, errorsBuffer),
	}
}
