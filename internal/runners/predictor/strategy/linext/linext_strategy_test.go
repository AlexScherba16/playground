package linext

import (
	tp "playground/internal/types"
	"reflect"
	"runtime"
	s "sync"
	"testing"
	"time"
)

type inputParameters struct {
	wg  *s.WaitGroup
	aCh tp.AggregatorChannel
	pCh tp.PredictorChannel
}

func TestLinearExtrapolationWorkerStrategy(t *testing.T) {
	/* ARRANGE */
	result := NewPredictWorkerStrategy()

	/* ACT */
	expected := reflect.ValueOf(linearExtrapolationWorker).Pointer()

	/* ASSERT */
	if reflect.ValueOf(result).Pointer() != expected {
		t.Fatalf("NewPredictWorkerStrategy() exp: %+v\ngot: %+v", expected, result)
	}
}

func TestLinearExtrapolationWorker_RunWorker(t *testing.T) {
	/* ARRANGE */
	aggrKey := "US"
	in := inputParameters{
		wg:  &s.WaitGroup{},
		aCh: tp.NewAggregatorChannel(0),
		pCh: tp.NewPredictorChannel(0),
	}
	defer close(in.pCh)
	// Prepare aggregated data
	aggregated := []*tp.AggregatedData{
		tp.NewAggregatedData(aggrKey, tp.LtvCollection{7, 0, 0, 0, 0, 0, 0}),
		tp.NewAggregatedData(aggrKey, tp.LtvCollection{1, 8, 0, 0, 0, 0, 0}),
		tp.NewAggregatedData(aggrKey, tp.LtvCollection{1, 4, 9, 12, 15, 0, 0}),
	}
	expected := tp.NewPredictedData(aggrKey, 180)
	in.wg.Add(1)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.aCh)
		for _, aggData := range aggregated {
			in.aCh <- aggData
		}
	}()
	go linearExtrapolationWorker(in.wg, aggrKey, in.aCh, in.pCh)

	/* ASSERT */
	for {
		select {
		// Assert expected predicted data
		case result, ok := <-in.pCh:
			if ok {
				// Assert result
				if result.Key() != expected.Key() {
					t.Fatalf("linearExtrapolationWorker() : expected key [%v], got [%v]", expected.Key(), result.Key())
				}
				if expected.Predicted() != result.Predicted() {
					t.Fatalf("linearExtrapolationWorker() : for %+v\nexpectd : %+v\ngot: %+v", result, expected.Predicted(), result.Predicted())
				}
				return
			} else {
				t.Fatalf("linearExtrapolationWorker() : shouldn't be here")
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestLinearExtrapolationWorker_RunWorkerWithCancelEvent(t *testing.T) {
	/* ARRANGE */
	aggrKey := "US"
	in := inputParameters{
		wg:  &s.WaitGroup{},
		aCh: tp.NewAggregatorChannel(0),
		pCh: tp.NewPredictorChannel(0),
	}
	defer close(in.pCh)
	// Prepare cansel event and store goroutines value before ACT stage
	aggregated := []*tp.AggregatedData{nil}
	expected := runtime.NumGoroutine()
	in.wg.Add(1)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.aCh)
		for _, aggData := range aggregated {
			in.aCh <- aggData
		}
	}()
	go linearExtrapolationWorker(in.wg, aggrKey, in.aCh, in.pCh)

	/* ASSERT */
	for {
		select {
		// Assert goroutines num
		case <-time.After(1 * time.Second):
			// linearExtrapolationWorker has no ability to notify about "Ok, I'm stopped"
			// It's better to reimplement worker, but time has pressure )
			// TODO: reimplement linearExtrapolationWorker, it should provide clear stopping notification

			result := runtime.NumGoroutine()
			if result != expected {
				t.Fatalf("linearExtrapolationWorker() goroutines value: expected : %+v\ngot: %+v",
					expected, result)
			}
			return
		}
	}
}
