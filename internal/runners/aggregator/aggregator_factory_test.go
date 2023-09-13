package aggregator

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
	"testing"
)

const (
	InvalidAggregateParameter = "aggregate_something"
)

func TestNewAggregator(t *testing.T) {
	tests := []struct {
		name          string
		aggregate     string
		expectedError bool
		errorStr      string
	}{
		{
			name:          "InvalidAggregatorParameter",
			aggregate:     InvalidAggregateParameter,
			expectedError: true,
			errorStr:      cerror.NewCustomError(fmt.Sprintf("%q invalid aggregate parameter", InvalidAggregateParameter)).Error(),
		},
		{
			name:      "AggregateCampaignParameter",
			aggregate: cnst.AggregateCampaign,
		},
		{
			name:      "AggregateCountryParameter",
			aggregate: cnst.AggregateCountry,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			// Prepare input parameters
			wg := &sync.WaitGroup{}
			recordCh := types.NewRecordChannel(0)
			aggregateCh := types.NewAggregateChannel(0)

			/* ACT */
			_, err := NewAggregator(wg, testCase.aggregate, recordCh, aggregateCh)

			/* ASSERT */
			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewAggregator() : expected error string [%s], got [%s]", testCase.errorStr, err.Error())
			}

			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewAggregator() : expected error %v, got %v", testCase.expectedError, err != nil)
			}
		})
	}
}
