package campaign

import (
	"fmt"
	cnst "playground/internal/constants"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"reflect"
	s "sync"
	"testing"
	"time"
)

const (
	InvalidCampaignAggregateParameter = "somethingThatDefinitelyNotCampaign"
)

type inputParameters struct {
	wg        *s.WaitGroup
	aggregate string
	rCh       tp.RecordChannel
	aCh       tp.AggregateChannel
}

type newAggregatorResult struct {
	aggregator *campaignAggregator
	err        error
}

func TestNewAggregator_InvalidInputParams(t *testing.T) {
	tests := []struct {
		name           string
		input          inputParameters
		expectedResult newAggregatorResult
		expectedError  bool
	}{
		{
			name:           "noWaitGroup",
			input:          inputParameters{nil, "", tp.NewRecordChannel(0), tp.NewAggregateChannel(0)},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid wait group")},
			expectedError:  true,
		},
		{
			name:  "invalidAggregateParam",
			input: inputParameters{&s.WaitGroup{}, InvalidCampaignAggregateParameter, tp.NewRecordChannel(0), tp.NewAggregateChannel(0)},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError(fmt.Sprintf("%q should be provided as an aggregate parameter, not %q",
				cnst.AggregateCampaign, InvalidCampaignAggregateParameter))},
			expectedError: true,
		},
		{
			name:           "noRecordChannel",
			input:          inputParameters{&s.WaitGroup{}, cnst.AggregateCampaign, nil, tp.NewAggregateChannel(0)},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid record channel")},
			expectedError:  true,
		},
		{
			name:           "noAggregateChannel",
			input:          inputParameters{&s.WaitGroup{}, cnst.AggregateCampaign, tp.NewRecordChannel(0), nil},
			expectedResult: newAggregatorResult{aggregator: nil, err: cerror.NewCustomError("invalid aggregate channel")},
			expectedError:  true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result, err := NewAggregator(testCase.input.wg, testCase.input.aggregate, testCase.input.rCh, testCase.input.aCh)

			/* ASSERT */
			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewAggregator() : expected error %v, got %v", testCase.expectedError, err != nil)
			}

			errorStr := ""
			if testCase.expectedResult.err != nil {
				errorStr = testCase.expectedResult.err.Error()
			}
			// Assert expected error string
			if (err != nil) && (err.Error() != errorStr) {
				t.Fatalf("NewAggregator() : expected error string [%s], got [%s]", errorStr, err.Error())
			}

			// Assert result
			if !reflect.DeepEqual(result, testCase.expectedResult.aggregator) {
				t.Fatalf("NewAggregator() exp: %+v\ngot: %+v", testCase.expectedResult.aggregator, result)
			}
		})
	}
}

func TestNewAggregator_ValidInputParams(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		cnst.AggregateCampaign,
		tp.NewRecordChannel(0),
		tp.NewAggregateChannel(0),
	}

	/* ACT */
	result, err := NewAggregator(in.wg, in.aggregate, in.rCh, in.aCh)
	// Assert unexpected error
	if err != nil {
		t.Fatalf("NewAggregator() : expected error string [%v], got [%v]", nil, err)
	}

	// Assert expected aggregator
	if result == nil {
		t.Fatalf("NewAggregator() : expected aggregator [%v], got [%v]", result, nil)
	}
}

func TestNewAggregator_RunReadRecordsChannel(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		cnst.AggregateCampaign,
		tp.NewRecordChannel(0),
		tp.NewAggregateChannel(0),
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
	aggregator, _ := NewAggregator(in.wg, in.aggregate, in.rCh, in.aCh)

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

func TestNewAggregator_RunReadCancelEventFromRecordsChannel(t *testing.T) {
	/* ARRANGE */
	in := inputParameters{
		&s.WaitGroup{},
		cnst.AggregateCampaign,
		tp.NewRecordChannel(0),
		tp.NewAggregateChannel(0),
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
	aggregator, _ := NewAggregator(in.wg, in.aggregate, in.rCh, in.aCh)

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
