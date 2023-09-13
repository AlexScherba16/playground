package main

import (
	"context"
	"log"
	"playground/internal/cli"
	cnst "playground/internal/constants"
	"playground/internal/runners/aggregator"
	"playground/internal/runners/datasource"
	"playground/internal/types"
	"sync"
)

func main() {
	// Get parsed user input flags
	flags, err := cli.NewFlags()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create channels storage
	ch := types.NewChannels(
		cnst.RecordChannelBuffer,
		cnst.ErrorChannelBuffer,
		cnst.AggregateChannelBuffer,
	)

	// Prepare input params
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Create datasource runner (Pipeline entry point)
	sourceRunner, err := datasource.NewDataSource(ctx, wg, flags.Source(), ch.RecordCh, ch.ErrorCh)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create aggregator runner
	aggregatorRunner, err := aggregator.NewAggregator(wg, flags.Aggregate(), ch.RecordCh, ch.AggregateCh)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Suppress unused variables, only for now, I promise )
	_ = cancel
	_ = sourceRunner
	_ = aggregatorRunner
}
