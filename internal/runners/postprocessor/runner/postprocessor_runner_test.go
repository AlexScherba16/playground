package runner

import (
	"fmt"
	"playground/internal/runners/postprocessor/strategy/campaign"
	"playground/internal/runners/postprocessor/strategy/country"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"reflect"
	"runtime"
	s "sync"
	"testing"
	"time"
)

type inputParameters struct {
	wg     *s.WaitGroup
	pCh    tp.PredictorChannel
	postCh tp.PostProcessorChannel
	pSt    tp.PostProcessorStrategy
}

type newPostProcessorResult struct {
	postProcessor *postProcessorRunner
	err           error
}

func TestNewPostProcessorRunner_InvalidInputParams(t *testing.T) {
	tests := []struct {
		name           string
		input          inputParameters
		expectedResult newPostProcessorResult
		expectedError  bool
	}{
		{
			name:           "noWaitGroup",
			input:          inputParameters{nil, nil, nil, nil},
			expectedResult: newPostProcessorResult{postProcessor: nil, err: cerror.NewCustomError("invalid wait group")},
			expectedError:  true,
		},
		{
			name:           "noPredictChannel",
			input:          inputParameters{&s.WaitGroup{}, nil, nil, nil},
			expectedResult: newPostProcessorResult{postProcessor: nil, err: cerror.NewCustomError("invalid predictor channel")},
			expectedError:  true,
		},
		{
			name:           "noPostProcessorChannel",
			input:          inputParameters{&s.WaitGroup{}, tp.NewPredictorChannel(0), nil, nil},
			expectedResult: newPostProcessorResult{postProcessor: nil, err: cerror.NewCustomError("invalid postprocessor channel")},
			expectedError:  true,
		},
		{
			name:           "noPredictStrategy",
			input:          inputParameters{&s.WaitGroup{}, tp.NewPredictorChannel(0), tp.NewPostProcessorChannel(0), nil},
			expectedResult: newPostProcessorResult{postProcessor: nil, err: cerror.NewCustomError("invalid postprocessor strategy")},
			expectedError:  true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result, err := NewPostProcessorRunner(testCase.input.wg, testCase.input.pCh, testCase.input.postCh, testCase.input.pSt)

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
			if !reflect.DeepEqual(result, testCase.expectedResult.postProcessor) {
				t.Fatalf("NewPredictor() exp: %+v\ngot: %+v", testCase.expectedResult.postProcessor, result)
			}
		})
	}
}

func TestNewPostProcessorRunner_ValidInputParamsCountryStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewPredictorChannel(0),
		tp.NewPostProcessorChannel(0),
		country.NewPostProcessorStrategy(),
	}

	/* ACT */
	result, err := NewPostProcessorRunner(in.wg, in.pCh, in.postCh, in.pSt)
	// Assert unexpected error
	if err != nil {
		t.Fatalf("NewPostProcessorRunner() : expected error string [%v], got [%v]", nil, err)
	}

	// Assert expected postProcessor
	if result == nil {
		t.Fatalf("NewPostProcessorRunner() : expected postProcessor [%v], got [%v]", result, nil)
	}
}

func TestNewPostProcessorRunner_RunWithCountryStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewPredictorChannel(0),
		tp.NewPostProcessorChannel(0),
		country.NewPostProcessorStrategy(),
	}
	//Prepare predicted data
	predicted := []*tp.PredictedData{
		tp.NewPredictedData("JP", 123.123),
		tp.NewPredictedData("US", 9999.99999),
	}
	// Iterate in reverse order cuz postprocessor sort data
	expectedPostProcData := []string{}
	for i := len(predicted) - 1; i >= 0; i-- {
		expectedPostProcData = append(expectedPostProcData,
			fmt.Sprintf("%s: %.2f", predicted[i].Key(), predicted[i].Predicted()))
	}

	in.wg.Add(1)
	postProcessor, _ := NewPostProcessorRunner(in.wg, in.pCh, in.postCh, in.pSt)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.pCh)
		for _, predictedData := range predicted {
			in.pCh <- predictedData
		}
	}()
	go postProcessor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected postprocessed data
		case result, ok := <-in.postCh:
			if ok {
				// Assert result
				expected := expectedPostProcData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedPostProcData = expectedPostProcData[1:]

			} else {
				// Assert empty expected records list
				if len(expectedPostProcData) != 0 {
					t.Fatalf("Run() unexpected postprocessor slice len exp: %+v\ngot: %+v", 0, len(expectedPostProcData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewPostProcessorRunner_RunWithCountryStrategyAndCancelEvent(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewPredictorChannel(0),
		tp.NewPostProcessorChannel(0),
		country.NewPostProcessorStrategy(),
	}
	//Prepare cancel event
	predicted := []*tp.PredictedData{nil}

	expectedGoroutines := runtime.NumGoroutine()
	in.wg.Add(1)
	postProcessor, _ := NewPostProcessorRunner(in.wg, in.pCh, in.postCh, in.pSt)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.pCh)
		for _, predictedData := range predicted {
			in.pCh <- predictedData
		}
	}()
	go postProcessor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected postprocessed data
		case _, ok := <-in.postCh:
			if ok {
				t.Fatalf("Run() shouldn't be here")
			} else {
				in.postCh = nil
			}
			// Assert potential hang situation
		// Assert goroutines num
		case <-time.After(15 * time.Millisecond):
			// postProcessor has no ability to notify about "Ok, I'm stopped"
			// TODO: reimplement postProcessor, it should provide clear stopping notification

			resultGoroutines := runtime.NumGoroutine()
			if resultGoroutines != expectedGoroutines {
				t.Fatalf("Run() goroutines value: expected : %+v\ngot: %+v",
					expectedGoroutines, resultGoroutines)
			}
			return
		}
	}
}

func TestNewPostProcessorRunner_RunWithCampaignStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewPredictorChannel(0),
		tp.NewPostProcessorChannel(0),
		campaign.NewPostProcessorStrategy(),
	}
	//Prepare predicted data
	predicted := []*tp.PredictedData{
		tp.NewPredictedData("JP", 123.123),
		tp.NewPredictedData("US", 9999.99999),
	}
	// Iterate in reverse order cuz postprocessor sort data
	expectedPostProcData := []string{}
	for i := len(predicted) - 1; i >= 0; i-- {
		expectedPostProcData = append(expectedPostProcData,
			fmt.Sprintf("<%s>: %.2f", predicted[i].Key(), predicted[i].Predicted()))
	}

	in.wg.Add(1)
	postProcessor, _ := NewPostProcessorRunner(in.wg, in.pCh, in.postCh, in.pSt)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.pCh)
		for _, predictedData := range predicted {
			in.pCh <- predictedData
		}
	}()
	go postProcessor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected postprocessed data
		case result, ok := <-in.postCh:
			if ok {
				// Assert result
				expected := expectedPostProcData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedPostProcData = expectedPostProcData[1:]

			} else {
				// Assert empty expected records list
				if len(expectedPostProcData) != 0 {
					t.Fatalf("Run() unexpected postprocessor slice len exp: %+v\ngot: %+v", 0, len(expectedPostProcData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewPostProcessorRunner_RunWithCampaignStrategyAndCancelEvent(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewPredictorChannel(0),
		tp.NewPostProcessorChannel(0),
		campaign.NewPostProcessorStrategy(),
	}
	// Prepare cancel event
	predicted := []*tp.PredictedData{nil}

	expectedGoroutines := runtime.NumGoroutine()
	in.wg.Add(1)
	postProcessor, _ := NewPostProcessorRunner(in.wg, in.pCh, in.postCh, in.pSt)

	/* ACT */
	// Mock aggregated streamer
	go func() {
		defer close(in.pCh)
		for _, predictedData := range predicted {
			in.pCh <- predictedData
		}
	}()
	go postProcessor.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected postprocessed data
		case _, ok := <-in.postCh:
			if ok {
				t.Fatalf("Run() shouldn't be here")
			} else {
				in.postCh = nil
			}
			// Assert potential hang situation
		// Assert goroutines num
		case <-time.After(15 * time.Millisecond):
			// postProcessor has no ability to notify about "Ok, I'm stopped"
			// TODO: reimplement postProcessor, it should provide clear stopping notification

			resultGoroutines := runtime.NumGoroutine()
			if resultGoroutines != expectedGoroutines {
				t.Fatalf("Run() goroutines value: expected : %+v\ngot: %+v",
					expectedGoroutines, resultGoroutines)
			}
			return
		}
	}
}
