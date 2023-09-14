package runner

import (
	"playground/internal/runners/aggregator/strategy/campaign_aggregator_strategy"
	"playground/internal/runners/aggregator/strategy/country_aggregator_strategy"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"reflect"
	s "sync"
	"testing"
	"time"
)

const (
	InvalidCountryAggregateParameter = "somethingThatDefinitelyNotCountry"
)

type inputParameters struct {
	wg       *s.WaitGroup
	rCh      tp.RecordChannel
	aCh      tp.AggregatorChannel
	strategy tp.AggregatorStrategy
}

type newAggregatorResult struct {
	aggregator *countryAggregator
	err        error
}

func TestNewAggregatorRunner_InvalidInputParams(t *testing.T) {
	tests := []struct {
		name           string
		input          inputParameters
		expectedResult newAggregatorResult
		expectedError  bool
	}{
		{
			name:           "noWaitGroup",
			input:          inputParameters{nil, nil, nil, nil},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid wait group")},
			expectedError:  true,
		},
		{
			name:           "noRecordChannel",
			input:          inputParameters{&s.WaitGroup{}, nil, nil, nil},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid record channel")},
			expectedError:  true,
		},
		{
			name:           "noAggregateChannel",
			input:          inputParameters{&s.WaitGroup{}, tp.NewRecordChannel(0), nil, nil},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid aggregate channel")},
			expectedError:  true,
		},
		{
			name:           "noAggregateStrategy",
			input:          inputParameters{&s.WaitGroup{}, tp.NewRecordChannel(0), tp.NewAggregatorChannel(0), nil},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid aggregation strategy")},
			expectedError:  true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result, err := NewAggregatorRunner(testCase.input.wg, testCase.input.rCh, testCase.input.aCh, testCase.input.strategy)

			/* ASSERT */
			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewAggregatorRunner() : expected error %v, got %v", testCase.expectedError, err != nil)
			}

			errorStr := ""
			if testCase.expectedResult.err != nil {
				errorStr = testCase.expectedResult.err.Error()
			}
			// Assert expected error string
			if (err != nil) && (err.Error() != errorStr) {
				t.Fatalf("NewAggregatorRunner() : expected error string [%s], got [%s]", errorStr, err.Error())
			}

			// Assert result
			if !reflect.DeepEqual(result, testCase.expectedResult.aggregator) {
				t.Fatalf("NewAggregatorRunner() exp: %+v\ngot: %+v", testCase.expectedResult.aggregator, result)
			}
		})
	}
}

func TestNewAggregatorRunner_ValidInputParamsCountryStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewRecordChannel(0),
		tp.NewAggregatorChannel(0),
		country_aggregator_strategy.NewCountryAggregatorStrategy(),
	}

	/* ACT */
	result, err := NewAggregatorRunner(in.wg, in.rCh, in.aCh, in.strategy)
	// Assert unexpected error
	if err != nil {
		t.Fatalf("NewAggregatorRunner() : expected error string [%v], got [%v]", nil, err)
	}

	// Assert expected aggregator
	if result == nil {
		t.Fatalf("NewAggregatorRunner() : expected aggregator [%v], got [%v]", result, nil)
	}
}

func TestNewAggregatorRunner_RunWithCountryStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewRecordChannel(0),
		tp.NewAggregatorChannel(0),
		country_aggregator_strategy.NewCountryAggregatorStrategy(),
	}
	// Prepare records and expected aggregated data
	records := []*tp.Record{
		tp.NewRecord("9566c74d-1003-4c4d-bbbb-0407d1e2c649", "JP", tp.LtvCollection{1.73305638789404, 1.7684248856061633, 2.781764692566589, 0, 0, 0, 0}),
		tp.NewRecord("6325253f-ec73-4dd7-a9e2-8bf921119c16", "US", tp.LtvCollection{1.9466884664338124, 3.166483202629052, 4.892883942338033, 0, 0, 0, 0}),
		tp.NewRecord("680b4e7c-8b76-4a1b-9d49-d4955c848621", "DE", tp.LtvCollection{1.281468676817884, 1.5047392622480078, 1.7456670792496436, 0, 0, 0, 0}),
	}

	expectedAggregatedData := []*tp.AggregatedData{}
	for _, record := range records {
		agg := tp.NewAggregatedData(record.Country(), record.Ltv())
		expectedAggregatedData = append(expectedAggregatedData, agg)
	}

	in.wg.Add(1)
	aggregator, _ := NewAggregatorRunner(in.wg, in.rCh, in.aCh, in.strategy)

	/* ACT */
	// Mock record streamer
	go func() {
		defer close(in.rCh)
		for _, record := range records {
			in.rCh <- record
		}
	}()
	go aggregator.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected aggregated data
		case result, ok := <-in.aCh:
			if ok {
				// Assert result
				expected := expectedAggregatedData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedAggregatedData = expectedAggregatedData[1:]

			} else {
				// Assert empty expected aggregated data list
				if len(expectedAggregatedData) != 0 {
					t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewAggregatorRunner_RunWithCountryStrategyAndCancelEvent(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewRecordChannel(0),
		tp.NewAggregatorChannel(0),
		country_aggregator_strategy.NewCountryAggregatorStrategy(),
	}
	// Prepare records and expected aggregated data
	records := []*tp.Record{
		tp.NewRecord("9566c74d-1003-4c4d-bbbb-0407d1e2c649", "JP", tp.LtvCollection{1.73305638789404, 1.7684248856061633, 2.781764692566589, 0, 0, 0, 0}),
		tp.NewRecord("6325253f-ec73-4dd7-a9e2-8bf921119c16", "US", tp.LtvCollection{1.9466884664338124, 3.166483202629052, 4.892883942338033, 0, 0, 0, 0}),
		nil,
	}

	expectedAggregatedData := []*tp.AggregatedData{}
	for _, record := range records[0 : len(records)-1] {
		agg := tp.NewAggregatedData(record.Country(), record.Ltv())
		expectedAggregatedData = append(expectedAggregatedData, agg)
	}

	in.wg.Add(1)
	aggregator, _ := NewAggregatorRunner(in.wg, in.rCh, in.aCh, in.strategy)

	/* ACT */
	// Mock record streamer
	go func() {
		defer close(in.rCh)
		for _, record := range records {
			in.rCh <- record
		}
	}()
	go aggregator.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected aggregated data
		case result, ok := <-in.aCh:
			if ok {
				// Assert cancel event
				if result == nil {
					if len(expectedAggregatedData) != 0 {
						t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
					}
					return
				}

				// Assert result
				expected := expectedAggregatedData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedAggregatedData = expectedAggregatedData[1:]

			} else {
				// Assert empty expected aggregated data list
				if len(expectedAggregatedData) != 0 {
					t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewAggregatorRunner_RunWithCampaignStrategy(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewRecordChannel(0),
		tp.NewAggregatorChannel(0),
		campaign_aggregator_strategy.NewCampaignAggregatorStrategy(),
	}
	// Prepare records and expected aggregated data
	records := []*tp.Record{
		tp.NewRecord("9566c74d-1003-4c4d-bbbb-0407d1e2c649", "JP", tp.LtvCollection{1.73305638789404, 1.7684248856061633, 2.781764692566589, 0, 0, 0, 0}),
		tp.NewRecord("6325253f-ec73-4dd7-a9e2-8bf921119c16", "US", tp.LtvCollection{1.9466884664338124, 3.166483202629052, 4.892883942338033, 0, 0, 0, 0}),
		tp.NewRecord("680b4e7c-8b76-4a1b-9d49-d4955c848621", "DE", tp.LtvCollection{1.281468676817884, 1.5047392622480078, 1.7456670792496436, 0, 0, 0, 0}),
	}

	expectedAggregatedData := []*tp.AggregatedData{}
	for _, record := range records {
		agg := tp.NewAggregatedData(record.CampaignId(), record.Ltv())
		expectedAggregatedData = append(expectedAggregatedData, agg)
	}

	in.wg.Add(1)
	aggregator, _ := NewAggregatorRunner(in.wg, in.rCh, in.aCh, in.strategy)

	/* ACT */
	// Mock record streamer
	go func() {
		defer close(in.rCh)
		for _, record := range records {
			in.rCh <- record
		}
	}()
	go aggregator.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected aggregated data
		case result, ok := <-in.aCh:
			if ok {
				// Assert result
				expected := expectedAggregatedData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedAggregatedData = expectedAggregatedData[1:]

			} else {
				// Assert empty expected aggregated data list
				if len(expectedAggregatedData) != 0 {
					t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewAggregatorRunner_RunWithCampaignStrategyAndCancelEvent(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		tp.NewRecordChannel(0),
		tp.NewAggregatorChannel(0),
		campaign_aggregator_strategy.NewCampaignAggregatorStrategy(),
	}
	// Prepare records and expected aggregated data
	records := []*tp.Record{
		tp.NewRecord("9566c74d-1003-4c4d-bbbb-0407d1e2c649", "JP", tp.LtvCollection{1.73305638789404, 1.7684248856061633, 2.781764692566589, 0, 0, 0, 0}),
		tp.NewRecord("6325253f-ec73-4dd7-a9e2-8bf921119c16", "US", tp.LtvCollection{1.9466884664338124, 3.166483202629052, 4.892883942338033, 0, 0, 0, 0}),
		nil,
	}

	expectedAggregatedData := []*tp.AggregatedData{}
	for _, record := range records[0 : len(records)-1] {
		agg := tp.NewAggregatedData(record.CampaignId(), record.Ltv())
		expectedAggregatedData = append(expectedAggregatedData, agg)
	}

	in.wg.Add(1)
	aggregator, _ := NewAggregatorRunner(in.wg, in.rCh, in.aCh, in.strategy)

	/* ACT */
	// Mock record streamer
	go func() {
		defer close(in.rCh)
		for _, record := range records {
			in.rCh <- record
		}
	}()
	go aggregator.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected aggregated data
		case result, ok := <-in.aCh:
			if ok {
				// Assert cancel event
				if result == nil {
					if len(expectedAggregatedData) != 0 {
						t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
					}
					return
				}

				// Assert result
				expected := expectedAggregatedData[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedAggregatedData = expectedAggregatedData[1:]

			} else {
				// Assert empty expected aggregated data list
				if len(expectedAggregatedData) != 0 {
					t.Fatalf("Run() unexpected aggregated data slice len exp: %+v\ngot: %+v", 0, len(expectedAggregatedData))
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}
