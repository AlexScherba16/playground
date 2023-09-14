package campaign_postprocessor_strategy

import (
	"fmt"
	t "playground/internal/types"
)

// campaignPostProcessor campaign_postprocessor_strategy predicted data conversion strategy function
func campaignPostProcessor(data *t.PredictedData) string {
	return fmt.Sprintf("<%s>: %.2f", data.Key(), data.Predicted())
}

// NewPostProcessorStrategy returns campaign_postprocessor_strategy predicted data convertor strategy
func NewPostProcessorStrategy() t.PostProcessorStrategy {
	return campaignPostProcessor
}
