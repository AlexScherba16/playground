package predictor

import (
	"math"
	"testing"
)

const Accuracy = 1e-9

func TestLinearExtrapolation(t *testing.T) {
	tests := []struct {
		data     []float64
		day      float64
		expected float64
	}{
		{
			data:     []float64{-2, -1, 0, 1, 2, 3, 4, 5},
			day:      6,
			expected: 3,
		},
		{
			data:     []float64{1, 2, 3, 4, 5, 6, 7},
			day:      100,
			expected: 100,
		},
		// Add here new cases, main idea was make sure that LinearExtrapolation is really linear )
	}

	for _, testCase := range tests {
		result := LinearExtrapolation(testCase.data, testCase.day)
		if math.Abs(result-testCase.expected) > Accuracy {
			t.Errorf("LinearExtrapolation() : input %v expected %v got %v", testCase.data, testCase.expected, result)
		}
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		data     []float64
		day      float64
		expected float64
	}{
		{
			data:     []float64{10, 10.6, 11.11, 15.91},
			day:      60,
			expected: 100,
		},
		// Add here new cases, main idea was make sure that Average is really average )
	}

	for _, testCase := range tests {
		result := Average(testCase.data, testCase.day)
		if math.Abs(result-testCase.expected) > 0.2 {
			t.Errorf("Average() : input %v expected %v got %v", testCase.data, testCase.expected, result)
		}
	}
}
