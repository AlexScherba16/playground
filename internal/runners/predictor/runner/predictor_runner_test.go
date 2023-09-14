package runner

import (
	"playground/internal/runners/predictor/strategy/linext"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"reflect"
	s "sync"
	"testing"
	"time"
)

type inputParameters struct {
	wg  *s.WaitGroup
	aCh tp.AggregatorChannel
	pCh tp.PredictorChannel
	pSt tp.PredictWorkerStrategy
}

type newPredictorResult struct {
	predictor *predictorRunner
	err       error
}

func TestNewPredictorRunner_InvalidInputParams(t *testing.T) {
	tests := []struct {
		name           string
		input          inputParameters
		expectedResult newPredictorResult
		expectedError  bool
	}{
		{
			name:           "noWaitGroup",
			input:          inputParameters{nil, nil, nil, nil},
			expectedResult: newPredictorResult{predictor: nil, err: cerror.NewCustomError("invalid wait group")},
			expectedError:  true,
		},
		{
			name:           "noAggregateChannel",
			input:          inputParameters{&s.WaitGroup{}, nil, nil, nil},
			expectedResult: newPredictorResult{predictor: nil, err: cerror.NewCustomError("invalid aggregator channel")},
			expectedError:  true,
		},
		{
			name:           "noPredictChannel",
			input:          inputParameters{&s.WaitGroup{}, tp.NewAggregatorChannel(0), nil, nil},
			expectedResult: newPredictorResult{predictor: nil, err: cerror.NewCustomError("invalid predictor channel")},
			expectedError:  true,
		},
		{
			name:           "noPredictStrategy",
			input:          inputParameters{&s.WaitGroup{}, tp.NewAggregatorChannel(0), tp.NewPredictorChannel(0), nil},
			expectedResult: newPredictorResult{predictor: nil, err: cerror.NewCustomError("invalid predictor strategy worker")},
			expectedError:  true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result, err := NewPredictorRunner(testCase.input.wg, testCase.input.aCh, testCase.input.pCh, testCase.input.pSt)

			/* ASSERT */
			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewPredictor() : expected error %v, got %v", testCase.expectedError, err != nil)
			}

			errorStr := ""
			if testCase.expectedResult.err != nil {
				errorStr = testCase.expectedResult.err.Error()
			}
			// Assert expected error string
			if (err != nil) && (err.Error() != errorStr) {
				t.Fatalf("NewPredictor() : expected error string [%s], got [%s]", errorStr, err.Error())
			}

			// Assert result
			if !reflect.DeepEqual(result, testCase.expectedResult.predictor) {
				t.Fatalf("NewPredictor() exp: %+v\ngot: %+v", testCase.expectedResult.predictor, result)
			}
		})
	}
}

func TestNewPredictorRunner_ValidInputParamsLinextWorkerStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewAggregatorChannel(0),
		tp.NewPredictorChannel(0),
		linext.NewPredictWorkerStrategy(),
	}

	/* ACT */
	result, err := NewPredictorRunner(in.wg, in.aCh, in.pCh, in.pSt)
	// Assert unexpected error
	if err != nil {
		t.Fatalf("NewPredictor() : expected error string [%v], got [%v]", nil, err)
	}

	// Assert expected predictor
	if result == nil {
		t.Fatalf("NewPredictor() : expected predictor [%v], got [%v]", result, nil)
	}
}

func TestNewPredictorRunner_RunWithLinextWorkerStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewAggregatorChannel(0),
		tp.NewPredictorChannel(0),
		linext.NewPredictWorkerStrategy(),
	}
	// Prepare aggregated data
	aggregated := []*tp.AggregatedData{
		tp.NewAggregatedData("JP", tp.LtvCollection{2, 4, 6, 8, 10, 0, 0}),
		tp.NewAggregatedData("US", tp.LtvCollection{3, 6, 9, 0, 0, 0, 0}),
		tp.NewAggregatedData("DE", tp.LtvCollection{1, 2, 3, 4, 5, 6, 7}),
	}
	expectedPredictedData := map[string]float64{
		"JP": 120, "US": 180, "DE": 60,
	}

	in.wg.Add(1)
	predictor, _ := NewPredictorRunner(in.wg, in.aCh, in.pCh, in.pSt)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.aCh)
		for _, aggData := range aggregated {
			in.aCh <- aggData
		}
	}()
	go predictor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected predicted data
		case result, ok := <-in.pCh:
			if ok {
				// Assert result
				value, found := expectedPredictedData[result.Key()]
				if !found {
					t.Fatalf("Run() no expected predicted data : %+v\ngot: %+v", expectedPredictedData, result.Key())
				}
				if value != result.Predicted() {
					t.Fatalf("Run() exp: %+v\ngot: %+v", value, result.Predicted())
				}
				// Delete expected data
				delete(expectedPredictedData, result.Key())

			} else {
				// Assert empty expected predicted map
				if len(expectedPredictedData) != 0 {
					t.Fatalf("Run() unexpected predicted data map len exp: %+v\ngot: %+v", 0, len(expectedPredictedData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewPredictorRunner_RunWithLinextWorkerStrategyAndCancelEvent(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewAggregatorChannel(0),
		tp.NewPredictorChannel(0),
		linext.NewPredictWorkerStrategy(),
	}
	// Prepare aggregated data and cancel event
	aggregated := []*tp.AggregatedData{
		tp.NewAggregatedData("JP", tp.LtvCollection{2, 4, 6, 8, 10, 0, 0}),
		tp.NewAggregatedData("US", tp.LtvCollection{3, 6, 9, 0, 0, 0, 0}),
		nil,
	}
	predictor, _ := NewPredictorRunner(in.wg, in.aCh, in.pCh, in.pSt)
	in.wg.Add(1)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.aCh)
		for _, aggData := range aggregated {
			in.aCh <- aggData
		}
	}()
	go predictor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected predicted data
		case result, ok := <-in.pCh:
			if ok {
				// Assert cansel event
				if result != nil {
					t.Fatalf("Run() exp: %+v\ngot: %+v", nil, result)
				}
				return

			} else {
				t.Fatalf("Run() : shouldn't be here")
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}
