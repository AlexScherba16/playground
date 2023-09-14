package campaign

import (
	"fmt"
	t "playground/internal/types"
)

// campaignPostProcessor campaign postprocessor strategy predicted data conversion strategy function
func campaignPostProcessor(data *t.PredictedData) string {
	return fmt.Sprintf("<%s>: %.2f", data.Key(), data.Predicted())
}

// NewPostProcessorStrategy returns campaign postprocessor strategy predicted data convertor strategy
func NewPostProcessorStrategy() t.PostProcessorStrategy {
	return campaignPostProcessor
}
