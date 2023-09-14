package predictor_factory

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
	"testing"
)

const (
	InvalidModelParameter = "PredictSomethingUnpredictable"
)

func TestNewRunner(t *testing.T) {
	tests := []struct {
		name          string
		model         string
		expectedError bool
		errorStr      string
	}{
		{
			name:          "InvalidModelParameter",
			model:         InvalidModelParameter,
			expectedError: true,
			errorStr:      cerror.NewCustomError(fmt.Sprintf("%q invalid model parameter", InvalidModelParameter)).Error(),
		},
		{
			name:  "LinearExtrapolationParameter",
			model: cnst.LinearExtrapolationPredictorModel,
		},
		{
			name:  "AverageParameter",
			model: cnst.AveragePredictorModel,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			// Prepare input parameters
			wg := &sync.WaitGroup{}
			aggregateCh := types.NewAggregatorChannel(0)
			predictCh := types.NewPredictorChannel(0)

			/* ACT */
			_, err := NewRunner(wg, testCase.model, aggregateCh, predictCh)

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
