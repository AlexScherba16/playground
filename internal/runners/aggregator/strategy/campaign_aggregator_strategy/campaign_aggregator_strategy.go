package campaign_aggregator_strategy

import (
	t "playground/internal/types"
)

// campaignAggregatorStrategy Record aggregation strategy, accorded to campaign key
func campaignAggregatorStrategy(record *t.Record) *t.AggregatedData {
	return t.NewAggregatedData(record.CampaignId(), record.Ltv())
}

// NewCampaignAggregatorStrategy returns campaign aggregator strategy
func NewCampaignAggregatorStrategy() t.AggregatorStrategy {
	return campaignAggregatorStrategy
}
