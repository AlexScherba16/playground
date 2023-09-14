package postprocessor_factory

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
	"testing"
)

const (
	InvalidPostProcessorParameter = "PostProcessSomethingUnPostProcessable"
)

func TestNewRunner(t *testing.T) {
	tests := []struct {
		name          string
		postProcessor string
		expectedError bool
		errorStr      string
	}{
		{
			name:          "InvalidPostProcessorParameter",
			postProcessor: InvalidPostProcessorParameter,
			expectedError: true,
			errorStr:      cerror.NewCustomError(fmt.Sprintf("%q invalid postprocessor parameter", InvalidPostProcessorParameter)).Error(),
		},
		{
			name:          "CountryPostProcessorParameter",
			postProcessor: cnst.AggregateCountry,
		},
		{
			name:          "CampaignPostProcessorParameter",
			postProcessor: cnst.AggregateCampaign,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			// Prepare input parameters
			wg := &sync.WaitGroup{}
			predictCh := types.NewPredictorChannel(0)
			postProcCh := types.NewPostProcessorChannel(0)

			/* ACT */
			_, err := NewRunner(wg, testCase.postProcessor, predictCh, postProcCh)

			/* ASSERT */
			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewRunner() : expected error string [%s], got [%s]", testCase.errorStr, err.Error())
			}

			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewRunner() : expected error %v, got %v", testCase.expectedError, err != nil)
			}
		})
	}
}
