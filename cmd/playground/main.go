package main

import (
	"context"
	"fmt"
	"log"
	"playground/internal/cli"
	cnst "playground/internal/constants"
	"playground/internal/runners/aggregator/aggregator_factory"
	"playground/internal/runners/common"
	"playground/internal/runners/datasource"
	postprocessor "playground/internal/runners/postprocessor/postprocessor_factory"
	"playground/internal/runners/predictor/predictor_factory"
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
		cnst.PredictChannelBuffer,
		cnst.PostProcessorChannelBuffer,
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
	aggregatorRunner, err := aggregator_factory.NewRunner(wg, flags.Aggregate(), ch.RecordCh, ch.AggregateCh)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create predictor runner
	predictorRunner, err := predictor_factory.NewRunner(wg, flags.Model(), ch.AggregateCh, ch.PredictCh)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create postprocessor runner
	postProcessorRunner, err := postprocessor.NewRunner(wg, flags.Aggregate(), ch.PredictCh, ch.PostProcCh)
	if err != nil {
		log.Fatalln(err.Error())
	}
	runners := []common.IRunner{
		sourceRunner,
		aggregatorRunner,
		predictorRunner,
		postProcessorRunner,
	}

	// Set wait group and launch runners
	wg.Add(len(runners))
	for _, runner := range runners {
		go runner.Run()
	}

	for {
		select {
		// Something went wrong, shutdown runners and show error message
		case err, ok := <-ch.ErrorCh:
			if ok {
				cancel()
				wg.Wait()
				log.Fatalln(err)
			} else {
				ch.ErrorCh = nil
			}
			// Read and print result
		case result, ok := <-ch.PostProcCh:
			if ok {
				fmt.Println(result)
			} else {
				return
			}

			// Simulate cancel event
			//case <-time.After(2 * time.Millisecond):
			//	log.Info("cancel )")
			//	cancel()
			//	wg.Wait()
			//	log.Info("application finished")
			//	return
		}
	}
}
