package average

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

func TestAverageWorkerStrategy(t *testing.T) {
	/* ARRANGE */
	result := NewPredictWorkerStrategy()

	/* ACT */
	expected := reflect.ValueOf(averageWorker).Pointer()

	/* ASSERT */
	if reflect.ValueOf(result).Pointer() != expected {
		t.Fatalf("NewPredictWorkerStrategy() exp: %+v\ngot: %+v", expected, result)
	}
}

func TestAverageWorker_RunWorker(t *testing.T) {
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
	expected := tp.NewPredictedData(aggrKey, 149.4)
	in.wg.Add(1)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.aCh)
		for _, aggData := range aggregated {
			in.aCh <- aggData
		}
	}()
	go averageWorker(in.wg, aggrKey, in.aCh, in.pCh)

	/* ASSERT */
	for {
		select {
		// Assert expected predicted data
		case result, ok := <-in.pCh:
			if ok {
				// Assert result
				if result.Key() != expected.Key() {
					t.Fatalf("averageWorker() : expected key [%v], got [%v]", expected.Key(), result.Key())
				}
				if expected.Predicted() != result.Predicted() {
					t.Fatalf("averageWorker() : for %+v\nexpectd : %+v\ngot: %+v", result, expected.Predicted(), result.Predicted())
				}
				return
			} else {
				t.Fatalf("averageWorker() : shouldn't be here")
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestAverageWorker_RunWorkerWithCancelEvent(t *testing.T) {
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
	go averageWorker(in.wg, aggrKey, in.aCh, in.pCh)

	/* ASSERT */
	for {
		select {
		// Assert goroutines num
		case <-time.After(1 * time.Second):
			// averageWorker has no ability to notify about "Ok, I'm stopped"
			// It's better to reimplement worker, but time has pressure )
			// TODO: reimplement averageWorker, it should provide clear stopping notification

			result := runtime.NumGoroutine()
			if result != expected {
				t.Fatalf("averageWorker() goroutines value: expected : %+v\ngot: %+v",
					expected, result)
			}
			return
		}
	}
}
