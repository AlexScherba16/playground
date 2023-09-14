package country_postprocessor_strategy

import (
	"fmt"
	t "playground/internal/types"
)

// countryPostProcessor country_postprocessor_strategy predicted data conversion strategy function
func countryPostProcessor(data *t.PredictedData) string {
	return fmt.Sprintf("%s: %.2f", data.Key(), data.Predicted())
}

// NewPostProcessorStrategy returns country_postprocessor_strategy predicted data convertor strategy
func NewPostProcessorStrategy() t.PostProcessorStrategy {
	return countryPostProcessor
}
